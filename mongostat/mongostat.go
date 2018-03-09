// Package mongostat provides an overview of the status of a currently running mongod or mongos instance.
package mongostat

import (
	"context"
	"mongo-tools/common/db"
	"mongo-tools/mongostat/stat_consumer"
	"mongo-tools/mongostat/stat_consumer/line"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"mongo-tools/common/options"

	"mongo-tools/mongostat/status"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ChannelClosed   int32 = 1
	ChannelUnclosed int32 = 0
)

// MongoStat is a container for the user-specified options and
// internal cluster state used for running mongostat.
type MongoStat struct {
	// Generic mongo tool options.
	Options *options.ToolOptions

	// Mongostat-specific output options.
	StatOptions *StatOptions

	// How long to sleep between printing the rows, and polling the server.
	SleepInterval time.Duration

	// New nodes can be "discovered" by any other node by sending a hostname
	// on this channel.
	//Discovered chan string

	// A map of hostname -> NodeMonitor for all the hosts that
	// are being monitored.
	Nodes map[string]*NodeMonitor

	// ClusterMonitor to manage collecting and printing the stats from all nodes.
	Cluster ClusterMonitor

	// Mutex to handle safe concurrent adding to or looping over discovered nodes.
	nodesLock sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMongoStat(ctx context.Context, opts *options.ToolOptions,
	sleep time.Duration, rowCount, during int64) (*MongoStat, error) {

	statOpts := &StatOptions{
		All:           true,
		HumanReadable: true,
		RowCount:      rowCount,
	}

	var factory stat_consumer.FormatterConstructor
	if statOpts.Json {
		factory = stat_consumer.FormatterConstructors["json"]
	} else {
		factory = stat_consumer.FormatterConstructors["grid"]
	}
	formatter := factory(statOpts.RowCount, !statOpts.NoHeaders)

	cliFlags := line.FlagHosts | line.FlagAlways

	if statOpts.All {
		cliFlags |= line.FlagAll
	}

	keyNames := line.DeprecatedKeyMap()

	readerConfig := &status.ReaderConfig{
		HumanReadable: statOpts.HumanReadable,
		TimeFormat:    "2006-01-02 15:04:05",
	}
	if statOpts.Json {
		readerConfig.TimeFormat = "15:04:05"
	}

	consumer := stat_consumer.NewStatConsumer(cliFlags, []string{}, keyNames, readerConfig, formatter)

	statctx, statcancel := context.WithCancel(ctx)

	var cluster ClusterMonitor
	if len(opts.Addrs) > 1 {
		cluster = &AsyncClusterMonitor{
			ReportChan:    make(chan *status.ServerStatus, len(opts.Addrs)),
			ErrorChan:     make(chan *status.NodeError, len(opts.Addrs)),
			LastStatLines: map[string]*line.StatLine{},
			Storage:       make(chan string, 10),
			StorageClosed: 0,
			Consumer:      consumer,
			startTime:     time.Now().Unix(),
			during:        during,
			ctx:           statctx,
		}
	} else {
		cluster = &SyncClusterMonitor{
			ReportChan:    make(chan *status.ServerStatus),
			ErrorChan:     make(chan *status.NodeError),
			Storage:       make(chan string, 10),
			StorageClosed: 0,
			Consumer:      consumer,
			startTime:     time.Now().Unix(),
			during:        during,
			ctx:           statctx,
		}
	}

	stat := &MongoStat{
		Options:       opts,
		StatOptions:   statOpts,
		Nodes:         map[string]*NodeMonitor{},
		SleepInterval: sleep,
		Cluster:       cluster,
		ctx:           statctx,
		cancel:        statcancel,
	}

	for _, v := range opts.Addrs {
		err := stat.AddNewNode(v)
		if err != nil {
			stat.cancel()
			return nil, err
		}
	}

	return stat, nil
}

// ConfigShard holds a mapping for the format of shard hosts as they
// appear in the config.shards collection.
type ConfigShard struct {
	Id   string `bson:"_id"`
	Host string `bson:"host"`
}

// NodeMonitor contains the connection pool for a single host and collects the
// mongostat data for that host on a regular interval.
type NodeMonitor struct {
	host, alias     string
	sessionProvider *db.SessionProvider

	// The time at which the node monitor last processed an update successfully.
	LastUpdate time.Time

	// The most recent error encountered when collecting stats for this node.
	Err error

	ctx context.Context
}

// SyncClusterMonitor is an implementation of ClusterMonitor that writes output
// synchronized with the timing of when the polling samples are collected.
// Only works with a single host at a time.
type SyncClusterMonitor struct {
	// Channel to listen for incoming stat data
	ReportChan chan *status.ServerStatus

	// Channel to listen for incoming errors
	ErrorChan chan *status.NodeError

	Storage       chan string
	StorageClosed int32
	lock          sync.Mutex

	// Creates and consumes StatLines using ServerStatuses
	Consumer *stat_consumer.StatConsumer

	//开始运行的时间
	startTime int64

	during int64

	ctx context.Context
}

// ClusterMonitor maintains an internal representation of a cluster's state,
// which can be refreshed with calls to Update(), and dumps output representing
// this internal state on an interval.
type ClusterMonitor interface {
	// Monitor() triggers monitoring and dumping output to begin
	// sleep is the interval to sleep between output dumps.
	// returns an error if it fails, and nil when monitoring ends
	Monitor(sleep time.Duration) error

	// Update signals the ClusterMonitor implementation to refresh its internal
	// state using the data contained in the provided ServerStatus.
	Update(stat *status.ServerStatus, err *status.NodeError)

	Message() (string, bool)

	Reset()
}

// AsyncClusterMonitor is an implementation of ClusterMonitor that writes output
// gotten from polling samples collected asynchronously from one or more servers.
type AsyncClusterMonitor struct {
	//Discover bool

	// Channel to listen for incoming stat data
	ReportChan chan *status.ServerStatus

	// Channel to listen for incoming errors
	ErrorChan chan *status.NodeError

	// Map of hostname -> latest stat data for the host
	LastStatLines map[string]*line.StatLine

	Storage       chan string
	StorageClosed int32

	// Mutex to protect access to LastStatLines
	mapLock sync.RWMutex

	// Creates and consumes StatLines using ServerStatuses
	Consumer *stat_consumer.StatConsumer

	//开始运行的时间
	startTime int64

	during int64

	ctx context.Context
}

func (cluster *SyncClusterMonitor) Message() (string, bool) {
	msg, ok := <-cluster.Storage
	return msg, ok
}

func (cluster *SyncClusterMonitor) Reset() {
	cluster.startTime = time.Now().Unix()
	cluster.Consumer.Reset()
}

func (cluster *SyncClusterMonitor) onceDone() {
	if atomic.LoadInt32(&cluster.StorageClosed) == ChannelClosed {
		return
	}

	cluster.lock.Lock()
	if cluster.StorageClosed == ChannelUnclosed {
		close(cluster.Storage)
		atomic.StoreInt32(&cluster.StorageClosed, ChannelClosed)
	}
	cluster.lock.Unlock()
}

// Update refreshes the internal state of the cluster monitor with the data
// in the StatLine. SyncClusterMonitor's implementation of Update blocks
// until it has written out its state, so that output is always dumped exactly
// once for each poll.
func (cluster *SyncClusterMonitor) Update(stat *status.ServerStatus, err *status.NodeError) {
	if err != nil {
		cluster.ErrorChan <- err
		return
	}
	cluster.ReportChan <- stat
}

// Monitor waits for data on the cluster's report channel. Once new data comes
// in, it formats and then displays it to stdout.
func (cluster *SyncClusterMonitor) Monitor(_ time.Duration) error {
	receivedData := false
	for {
		var statLine *line.StatLine
		var ok bool
		select {
		case stat := <-cluster.ReportChan:
			statLine, ok = cluster.Consumer.Update(stat)
			if !ok {
				continue
			}
		case err := <-cluster.ErrorChan:
			if !receivedData {
				statLine = &line.StatLine{
					Error:  err,
					Fields: map[string]string{"host": err.Host},
				}
				str, _ := cluster.Consumer.FormatLines([]*line.StatLine{statLine})
				cluster.Storage <- str
				//sleep 100ms 让channel有时间反应，将错误数据发送出去
				time.Sleep(100 * time.Millisecond)
				cluster.onceDone()
				return err
			}
			statLine = &line.StatLine{
				Error:  err,
				Fields: map[string]string{"host": err.Host},
			}
		case <-cluster.ctx.Done():
			cluster.onceDone()
			//fmt.Println("sync cluster monitor ctx done")
			return nil
		}
		receivedData = true
		str, finish := cluster.Consumer.FormatLines([]*line.StatLine{statLine})
		cluster.Storage <- str
		timeout := time.Now().Unix()-cluster.startTime > cluster.during
		if finish || timeout {
			cluster.onceDone()
			return nil
		}
	}
}

func (cluster *AsyncClusterMonitor) Message() (string, bool) {
	msg, ok := <-cluster.Storage
	return msg, ok
}

func (cluster *AsyncClusterMonitor) Reset() {
	cluster.startTime = time.Now().Unix()
	cluster.Consumer.Reset()
}

func (cluster *AsyncClusterMonitor) onceDone() {
	if atomic.LoadInt32(&cluster.StorageClosed) == ChannelClosed {
		return
	}

	cluster.mapLock.Lock()
	if cluster.StorageClosed == ChannelUnclosed {
		close(cluster.Storage)
		atomic.StoreInt32(&cluster.StorageClosed, ChannelClosed)
	}
	cluster.mapLock.Unlock()
}

// updateHostInfo updates the internal map with the given StatLine data.
// Safe for concurrent access.
func (cluster *AsyncClusterMonitor) updateHostInfo(stat *line.StatLine) {
	cluster.mapLock.Lock()
	defer cluster.mapLock.Unlock()
	host := stat.Fields["host"]
	cluster.LastStatLines[host] = stat
}

// printSnapshot formats and dumps the current state of all the stats collected.
// returns whether the program should now exit
func (cluster *AsyncClusterMonitor) printSnapshot() bool {
	cluster.mapLock.RLock()
	defer cluster.mapLock.RUnlock()
	lines := make([]*line.StatLine, 0, len(cluster.LastStatLines))
	for _, stat := range cluster.LastStatLines {
		lines = append(lines, stat)
	}
	if len(lines) == 0 {
		return false
	}

	str, finish := cluster.Consumer.FormatLines(lines)
	cluster.Storage <- str
	return finish
}

// Update sends a new StatLine on the cluster's report channel.
func (cluster *AsyncClusterMonitor) Update(stat *status.ServerStatus, err *status.NodeError) {
	if err != nil {
		cluster.ErrorChan <- err
		return
	}
	cluster.ReportChan <- stat
}

// The Async implementation of Monitor starts the goroutines that listen for incoming stat data,
// and dump snapshots at a regular interval.
func (cluster *AsyncClusterMonitor) Monitor(sleep time.Duration) error {
	select {
	case stat := <-cluster.ReportChan:
		cluster.Consumer.Update(stat)
	case err := <-cluster.ErrorChan:
		// error out if the first result is an error
		statLine := &line.StatLine{
			Error:  err,
			Fields: map[string]string{"host": err.Host},
		}
		cluster.updateHostInfo(statLine)

		n := len(cluster.ErrorChan)
		for i := 0; i < n; i++ {
			err := <-cluster.ErrorChan
			statLine := &line.StatLine{
				Error:  err,
				Fields: map[string]string{"host": err.Host},
			}
			cluster.updateHostInfo(statLine)
		}
		cluster.printSnapshot()
		//sleep 100ms 让channel有时间反应，将错误数据发送出去
		time.Sleep(100 * time.Millisecond)
		cluster.onceDone()
		return err
	}

	go func() {
		for {
			select {
			case stat := <-cluster.ReportChan:
				statLine, ok := cluster.Consumer.Update(stat)
				if ok {
					cluster.updateHostInfo(statLine)
				}
			case err := <-cluster.ErrorChan:
				statLine := &line.StatLine{
					Error:  err,
					Fields: map[string]string{"host": err.Host},
				}
				cluster.updateHostInfo(statLine)
			case <-cluster.ctx.Done():
				cluster.onceDone()
				//fmt.Println("async cluster monitor goroutine ctx done")
				return
			}
		}
	}()

	ticker := time.NewTicker(sleep)
	for {
		select {
		case <-ticker.C:
			//如果达到退出条件，输出了指定行数|达到指定的持续时间
			timeout := time.Now().Unix()-cluster.startTime > cluster.during
			if timeout || cluster.printSnapshot() {
				cluster.onceDone()
				return nil
			}
		case <-cluster.ctx.Done():
			cluster.onceDone()
			//fmt.Println("async cluster monitor ctx done")
			return nil
		}
	}
	return nil
}

// NewNodeMonitor copies the same connection settings from an instance of
// ToolOptions, but monitors fullHost.
func NewNodeMonitor(opts *options.ToolOptions, fullHost string) (*NodeMonitor, error) {
	optsCopy := options.New("mongostat")

	optsCopy.AppName = opts.AppName
	optsCopy.Source = opts.Source
	optsCopy.Username = opts.Username
	optsCopy.Password = opts.Password
	optsCopy.ReplicaSetName = opts.ReplicaSetName
	optsCopy.Timeout = opts.Timeout
	optsCopy.TCPKeepAliveSeconds = opts.TCPKeepAliveSeconds

	optsCopy.Addrs = []string{fullHost}
	//直连每个host
	optsCopy.Direct = true
	sessionProvider, err := db.NewSessionProvider(optsCopy)
	if err != nil {
		return nil, err
	}

	return &NodeMonitor{
		host:            fullHost,
		sessionProvider: sessionProvider,
		LastUpdate:      time.Now().Local(),
		Err:             nil,
	}, nil
}

// Report collects the stat info for a single node and sends found hostnames on
// the "discover" channel if checkShards is true.
func (node *NodeMonitor) Poll() (*status.ServerStatus, error) {
	stat := &status.ServerStatus{}
	s, err := node.sessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	// The read pref for the session must be set to 'secondary' to enable using
	// the driver with 'direct' connections, which disables the built-in
	// replset discovery mechanism since we do our own node discovery here.
	s.SetMode(mgo.Eventual, true)

	// Disable the socket timeout - otherwise if db.serverStatus() takes a long time on the server
	// side, the client will close the connection early and report an error.
	s.SetSocketTimeout(0)
	defer s.Close()

	err = s.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, stat)
	if err != nil {
		return nil, err
	}

	node.Err = nil
	stat.SampleTime = time.Now().Local()

	node.alias = stat.Host
	stat.Host = node.host

	return stat, nil
}

// Watch continuously collects and processes stats for a single node on a
// regular interval. At each interval, it triggers the node's Poll function
// with the 'discover' channel.
func (node *NodeMonitor) Watch(sleep time.Duration, cluster ClusterMonitor) {

	ticker := time.NewTicker(sleep)
	for {
		select {
		case <-ticker.C:
			stat, err := node.Poll()

			var nodeError *status.NodeError
			if err != nil {
				//fmt.Println("Poll error: ", err.Error())
				nodeError = status.NewNodeError(node.host, err)
			}
			cluster.Update(stat, nodeError)
		case <-node.ctx.Done():
			//fmt.Println(node.host, "ctx done")
			return
		}
	}
}

func parseHostPort(fullHostName string) (string, string) {
	if colon := strings.LastIndex(fullHostName, ":"); colon >= 0 {
		return fullHostName[0:colon], fullHostName[colon+1:]
	}
	return fullHostName, "27017"
}

// AddNewNode adds a new host name to be monitored and spawns the necessary
// goroutine to collect data from it.
func (mstat *MongoStat) AddNewNode(fullhost string) error {
	mstat.nodesLock.Lock()
	defer mstat.nodesLock.Unlock()

	// Remove the 'shardXX/' prefix from the hostname, if applicable
	pieces := strings.Split(fullhost, "/")
	fullhost = pieces[len(pieces)-1]

	if _, hasKey := mstat.Nodes[fullhost]; hasKey {
		return nil
	}
	for _, node := range mstat.Nodes {
		if node.alias == fullhost {
			return nil
		}
	}
	// Create a new node monitor for this host
	node, err := NewNodeMonitor(mstat.Options, fullhost)
	if err != nil {
		return err
	}

	node.ctx = mstat.ctx

	mstat.Nodes[fullhost] = node
	//go node.Watch(mstat.SleepInterval, mstat.Discovered, mstat.Cluster)
	go node.Watch(mstat.SleepInterval, mstat.Cluster)
	return nil
}

// Run is the top-level function that starts the monitoring
// and discovery goroutines
// https://docs.mongodb.com/v3.2/reference/program/mongostat/
func (mstat *MongoStat) Run() error {

	return mstat.Cluster.Monitor(mstat.SleepInterval)
}

func (mstat *MongoStat) Reset() {
	mstat.Cluster.Reset()
}

func (mstat *MongoStat) Stop() {
	if mstat.cancel != nil {
		mstat.cancel()
		mstat.cancel = nil
	}

	//给ctx一点响应时间
	time.Sleep(100 * time.Millisecond)

	for _, node := range mstat.Nodes {
		node.sessionProvider.Close()
	}
	//fmt.Println("Mongo Session Closed")
}
