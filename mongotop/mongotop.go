package mongotop

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"
)

const (
	ChannelClosed   int32 = 1
	ChannelUnclosed int32 = 0
)

// MongoTop is a container for the user-specified options and
// internal state used for running mongotop.
type MongoTop struct {
	// Generic mongo tool options
	Options *options.ToolOptions

	// Mongotop-specific output options
	OutputOptions *Output

	// for connecting to the db
	SessionProvider *db.SessionProvider

	// Length of time to sleep between each polling.
	Sleeptime time.Duration

	//持续时间
	During int64

	Storage       chan string
	StorageClosed int32
	lock          sync.Mutex

	numPrinted int32
	startTime  int64
	ctx        context.Context
	cancel     context.CancelFunc

	previousServerStatus *ServerStatus
	previousTop          *Top
}

//mongotop  3.0.0+ not support --locks
func NewMongoTop(ctx context.Context, opts *options.ToolOptions, oopts *Output, sp *db.SessionProvider,
	st time.Duration, during int64) *MongoTop {

	top := &MongoTop{
		Options:         opts,
		OutputOptions:   oopts,
		SessionProvider: sp,
		Sleeptime:       st,
		During:          during,
		Storage:         make(chan string, 10),
		StorageClosed:   0,
		numPrinted:      0,
		startTime:       time.Now().Unix(),
	}

	top.ctx, top.cancel = context.WithCancel(ctx)

	return top
}

func (mt *MongoTop) onceDone() {
	if atomic.LoadInt32(&mt.StorageClosed) == ChannelClosed {
		return
	}

	mt.lock.Lock()
	if mt.StorageClosed == ChannelUnclosed {
		close(mt.Storage)
		atomic.StoreInt32(&mt.StorageClosed, ChannelClosed)
	}
	mt.lock.Unlock()
}

func (mt *MongoTop) runDiff() (outDiff FormattableDiff, err error) {
	session, err := mt.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	var currentServerStatus ServerStatus
	var currentTop Top
	commandName := "top"
	var dest interface{} = &currentTop
	if mt.OutputOptions.Locks {
		commandName = "serverStatus"
		dest = &currentServerStatus
	}

	err = session.DB("admin").Run(commandName, dest)

	if err != nil {
		mt.previousServerStatus = nil
		mt.previousTop = nil
		return nil, err
	}

	if mt.OutputOptions.Locks {
		if currentServerStatus.Locks == nil {
			return nil, fmt.Errorf("server does not support reporting lock information")
		}
		for _, ns := range currentServerStatus.Locks {
			if ns.AcquireCount != nil {
				return nil, fmt.Errorf("server does not support reporting lock information")
			}
		}
		if mt.previousServerStatus != nil {
			serverStatusDiff := currentServerStatus.Diff(*mt.previousServerStatus)
			outDiff = serverStatusDiff
		}
		mt.previousServerStatus = &currentServerStatus
	} else {
		if mt.previousTop != nil {
			topDiff := currentTop.Diff(*mt.previousTop)
			outDiff = topDiff
		}
		mt.previousTop = &currentTop
	}
	return outDiff, nil
}

func (mt *MongoTop) IsFinished() bool {
	if mt.OutputOptions.RowCount > 0 && atomic.LoadInt32(&mt.numPrinted) > mt.OutputOptions.RowCount {
		return true
	}

	if time.Now().Unix()-mt.startTime > mt.During {
		return true
	}

	return false
}

// Run executes the mongotop program.
// https://docs.mongodb.com/v3.2/reference/program/mongotop/
// https://docs.mongodb.com/v3.2/reference/command/top/
func (mt *MongoTop) Run() error {

	hasData := false
	ticker := time.NewTicker(mt.Sleeptime)

	//释放mongo连接资源
	defer func() {
		//mt.SessionProvider.Close()
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if mt.IsFinished() {
				mt.onceDone()
				return nil
			}

			atomic.AddInt32(&mt.numPrinted, 1)

			diff, err := mt.runDiff()
			if err != nil {
				// If this is the first time trying to poll the server and it fails,
				// just stop now instead of trying over and over.
				if !hasData {
					mt.onceDone()
					return err
				}
			}

			hasData = true

			if diff != nil {
				if mt.OutputOptions.Json {
					mt.Storage <- diff.JSON()
				} else {
					mt.Storage <- diff.Grid()
				}
			}
		case <-mt.ctx.Done():
			mt.onceDone()
			//fmt.Println("mongotop ctx done")
			return nil
		}
	}
}

func (mt *MongoTop) Reset() {
	atomic.StoreInt32(&mt.numPrinted, 0)
	mt.startTime = time.Now().Unix()
}

func (mt *MongoTop) Stop() {
	if mt.cancel != nil {
		mt.cancel()
		mt.cancel = nil
	}
	mt.SessionProvider.Close()
}
