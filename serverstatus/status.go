package serverstatus

import "time"

type ServerStatusInfo struct {
	Host               string                 `bson:"host" json:"host"`
	Version            string                 `bson:"version" json:"version"`
	Process            string                 `bson:"process" json:"process"`
	AdvisoryHostFQDNs  []string               `bson:"advisoryHostFQDNs" json:"advisoryHostFQDNs"`
	Pid                int64                  `bson:"pid" json:"pid"`
	Uptime             int64                  `bson:"uptime" json:"uptime"`
	UptimeMillis       int64                  `bson:"uptimeMillis" json:"uptimeMillis"`
	UptimeEstimate     int64                  `bson:"uptimeEstimate" json:"uptimeEstimate"`
	LocalTime          time.Time              `bson:"localTime" json:"localTime"`
	Asserts            map[string]int64       `bson:"asserts" json:"asserts"`
	BackgroundFlushing *FlushStats            `bson:"backgroundFlushing" json:"backgroundFlushing"`
	ExtraInfo          *ExtraInfo             `bson:"extra_info" json:"extra_info"`
	Connections        *ConnectionStats       `bson:"connections" json:"connections"`
	Dur                *DurStats              `bson:"dur" json:"dur"`
	GlobalLock         *GlobalLockStats       `bson:"globalLock" json:"globalLock"`
	Locks              map[string]LockStats   `bson:"locks,omitempty" json:"locks,omitempty"`
	Network            *NetworkStats          `bson:"network" json:"network"`
	Opcounters         *OpcountStats          `bson:"opcounters" json:"opcounters"`
	OpcountersRepl     *OpcountStats          `bson:"opcountersRepl" json:"opcountersRepl"`
	RecordStats        *DBRecordStats         `bson:"recordStats" json:"recordStats"`
	Mem                *MemStats              `bson:"mem" json:"mem"`
	Repl               *ReplStatus            `bson:"repl" json:"repl"`
	ShardCursorType    map[string]interface{} `bson:"shardCursorType" json:"shardCursorType"`
	StorageEngine      *StorageEngine         `bson:"storageEngine" json:"storageEngine"`
	WiredTiger         *WiredTiger            `bson:"wiredTiger" json:"wiredTiger"`
}

type StorageEngine struct {
	Name                   string `bson:"name" json:"name"`
	Persistent             bool   `bson:"persistent" json:"persistent"`
	SupportsCommittedReads bool   `bson:"supportsCommittedReads" json:"supportsCommittedReads"`
}

// WiredTiger stores information related to the WiredTiger storage engine.
type WiredTiger struct {
	Transaction TransactionStats       `bson:"transaction" json:"transaction"`
	Concurrent  ConcurrentTransactions `bson:"concurrentTransactions" json:"concurrentTransactions"`
	Cache       CacheStats             `bson:"cache" json:"cache"`
	Async       Async                  `bson:"async" json:"async"`
}

type Async struct {
	CurrentQueueLength int `bson:"current work queue length" json:"current work queue length"`
	MaxQueueLength     int `bson:"maximum work queue length" json:"maximum work queue length"`
	WorkerNoWork       int `bson:"number of times worker found no work" json:"number of times worker found no work"`
	OpFailed           int `bson:"number of times operation allocation failed" json:"number of times operation allocation failed"`
	TotalAlloc         int `bson:"total allocations" json:"total allocations"`
	TotalCompact       int `bson:"total compact calls" json:"total compact calls"`
	TotalInsert        int `bson:"total insert calls" json:"total insert calls"`
	TotalRemove        int `bson:"total remove calls" json:"total remove calls"`
	TotalSearch        int `bson:"total search calls" json:"total search calls"`
	TotalUpdate        int `bson:"total update calls" json:"total update calls"`
}

type ConcurrentTransactions struct {
	Write ConcurrentTransStats `bson:"write" json:"write"`
	Read  ConcurrentTransStats `bson:"read" json:"read"`
}

type ConcurrentTransStats struct {
	Out          int64 `bson:"out" json:"out"`
	Current      int64 `bson:"current" json:"current"`
	Available    int64 `bson:"available" json:"available"`
	TotalTickets int64 `bson:"totalTickets" json:"totalTickets"`
}

// CacheStats stores cache statistics for WiredTiger.
type CacheStats struct {
	TrackedDirtyBytes  int64 `bson:"tracked dirty bytes in the cache" json:"tracked dirty bytes in the cache"`
	CurrentCachedBytes int64 `bson:"bytes currently in the cache" json:"bytes currently in the cache"`
	MaxBytesConfigured int64 `bson:"maximum bytes configured" json:"maximum bytes configured"`
}

// TransactionStats stores transaction checkpoints in WiredTiger.
type TransactionStats struct {
	TransCheckpoints int64 `bson:"transaction checkpoints" json:"transaction checkpoints"`
}

// ReplStatus stores data related to replica sets.
type ReplStatus struct {
	SetName    string   `bson:"setName" json:"setName"`
	IsMaster   bool     `bson:"ismaster" json:"ismaster"`
	Secondary  bool     `bson:"secondary" json:"secondary"`
	Primary    string   `bson:"primary" json:"primary"`
	SetVersion int      `bson:"setVersion" json:"setVersion"`
	Rbid       int64    `bson:"rbid" json:"rbid"`
	Hosts      []string `bson:"hosts" json:"hosts"`
	Me         string   `bson:"me" json:"me"`
}

// DBRecordStats stores data related to memory operations across databases.
type DBRecordStats struct {
	AccessesNotInMemory       int64                     `bson:"accessesNotInMemory" json:"accessesNotInMemory"`
	PageFaultExceptionsThrown int64                     `bson:"pageFaultExceptionsThrown" json:"pageFaultExceptionsThrown"`
	DBRecordAccesses          map[string]RecordAccesses `bson:",inline" json:"dBRecordAccesses"`
}

// RecordAccesses stores data related to memory operations scoped to a database.
type RecordAccesses struct {
	AccessesNotInMemory       int64 `bson:"accessesNotInMemory" json:"accessesNotInMemory"`
	PageFaultExceptionsThrown int64 `bson:"pageFaultExceptionsThrown" json:"pageFaultExceptionsThrown"`
}

// MemStats stores data related to memory statistics.
type MemStats struct {
	Bits              int64       `bson:"bits" json:"bits"`
	Resident          int64       `bson:"resident" json:"resident"`
	Virtual           int64       `bson:"virtual" json:"virtual"`
	Supported         interface{} `bson:"supported" json:"supported"`
	Mapped            int64       `bson:"mapped" json:"mapped"`
	MappedWithJournal int64       `bson:"mappedWithJournal" json:"mappedWithJournal"`
}

// FlushStats stores information about memory flushes.
type FlushStats struct {
	Flushes      int64     `bson:"flushes" json:"flushes"`
	TotalMs      int64     `bson:"total_ms" json:"total_ms"`
	AverageMs    float64   `bson:"average_ms" json:"average_ms"`
	LastMs       int64     `bson:"last_ms" json:"last_ms"`
	LastFinished time.Time `bson:"last_finished" json:"last_finished"`
}

// ConnectionStats stores information related to incoming database connections.
type ConnectionStats struct {
	Current      int64 `bson:"current" json:"current"`
	Available    int64 `bson:"available" json:"available"`
	TotalCreated int64 `bson:"totalCreated" json:"totalCreated"`
}

// DurTiming stores information related to journaling.
type DurTiming struct {
	Dt               int64 `bson:"dt" json:"dt"`
	PrepLogBuffer    int64 `bson:"prepLogBuffer" json:"prepLogBuffer"`
	WriteToJournal   int64 `bson:"writeToJournal" json:"writeToJournal"`
	WriteToDataFiles int64 `bson:"writeToDataFiles" json:"writeToDataFiles"`
	RemapPrivateView int64 `bson:"remapPrivateView" json:"remapPrivateView"`
}

// DurStats stores information related to journaling statistics.
type DurStats struct {
	Commits            int64     `bson:"commits" json:"commits"`
	JournaledMB        int64     `bson:"journaledMB" json:"journaledMB"`
	WriteToDataFilesMB int64     `bson:"writeToDataFilesMB" json:"writeToDataFilesMB"`
	Compression        int64     `bson:"compression" json:"compression"`
	CommitsInWriteLock int64     `bson:"commitsInWriteLock" json:"commitsInWriteLock"`
	EarlyCommits       int64     `bson:"earlyCommits" json:"earlyCommits"`
	TimeMs             DurTiming `json:"timeMs"`
}

// QueueStats stores the number of queued read/write operations.
type QueueStats struct {
	Total   int64 `bson:"total" json:"total"`
	Readers int64 `bson:"readers" json:"readers"`
	Writers int64 `bson:"writers" json:"writers"`
}

// ClientStats stores the number of active read/write operations.
type ClientStats struct {
	Total   int64 `bson:"total" json:"total"`
	Readers int64 `bson:"readers" json:"readers"`
	Writers int64 `bson:"writers" json:"writers"`
}

// GlobalLockStats stores information related locks in the MMAP storage engine.
type GlobalLockStats struct {
	TotalTime     int64        `bson:"totalTime" json:"totalTime"`
	LockTime      int64        `bson:"lockTime" json:"lockTime"`
	CurrentQueue  *QueueStats  `bson:"currentQueue" json:"currentQueue"`
	ActiveClients *ClientStats `bson:"activeClients" json:"activeClients"`
}

// NetworkStats stores information related to network traffic.
type NetworkStats struct {
	BytesIn     int64 `bson:"bytesIn" json:"bytesIn"`
	BytesOut    int64 `bson:"bytesOut" json:"bytesOut"`
	NumRequests int64 `bson:"numRequests" json:"numRequests"`
}

// OpcountStats stores information related to comamnds and basic CRUD operations.
type OpcountStats struct {
	Insert  int64 `bson:"insert" json:"insert"`
	Query   int64 `bson:"query" json:"query"`
	Update  int64 `bson:"update" json:"update"`
	Delete  int64 `bson:"delete" json:"delete"`
	GetMore int64 `bson:"getmore" json:"getmore"`
	Command int64 `bson:"command" json:"command"`
}

// ReadWriteLockTimes stores time spent holding read/write locks.
type ReadWriteLockTimes struct {
	Read       int64 `bson:"R" json:"R"`
	Write      int64 `bson:"W" json:"W"`
	ReadLower  int64 `bson:"r" json:"r"`
	WriteLower int64 `bson:"w" json:"w"`
}

// LockStats stores information related to time spent acquiring/holding locks
// for a given database.
type LockStats struct {
	DeadlockCount       ReadWriteLockTimes `bson:"deadlockCount" json:"deadlockCount"`
	TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros" json:"timeAcquiringMicros"`

	// AcquireCount and AcquireWaitCount are new fields of the lock stats only populated on 3.0 or newer.
	// Typed as a pointer so that if it is nil, mongostat can assume the field is not populated
	// with real namespace data.
	AcquireCount     *ReadWriteLockTimes `bson:"acquireCount,omitempty" json:"acquireCount,omitempty"`
	AcquireWaitCount *ReadWriteLockTimes `bson:"acquireWaitCount,omitempty" json:"acquireWaitCount,omitempty"`
}

// ExtraInfo stores additional platform specific information.
type ExtraInfo struct {
	Note            string `bson:"note" json:"note"`
	PageFaults      int64  `bson:"page_faults" json:"page_faults"`
	UsagePageFileMB int    `bson:"usagePageFileMB" json:"usagePageFileMB"`
	TotalPageFileMB int    `bson:"totalPageFileMB" json:"totalPageFileMB"`
	AvailPageFileMB int    `bson:"availPageFileMB" json:"availPageFileMB"`
	RamMB           int    `bson:"ramMB" json:"ramMB"`
}
