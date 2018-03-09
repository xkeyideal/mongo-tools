package collstat

import (
	"mongo-tools/common/db"
	"mongo-tools/common/options"

	"gopkg.in/mgo.v2/bson"
)

type CollectionStat struct {
	Ns             string      `bson:"ns" json:"ns"`
	Count          int         `bson:"count" json:"count"`
	Size           int         `bson:"size" json:"size"`
	AvgObjSize     int         `bson:"avgObjSize" json:"avgObjSize"`
	StorageSize    int         `bson:"storageSize" json:"storageSize"`
	Capped         bool        `bson:"capped" json:"capped"`
	Max            int         `bson:"max" json:"max"`
	MaxSize        int         `bson:"maxSize" json:"maxSize"`
	WiredTiger     *WiredTiger `bson:"wiredTiger" json:"wiredTiger"`
	Nindexes       int         `bson:"nindexes" json:"nindexes"`
	totalIndexSize int64       `bson:"totalIndexSize" json:"totalIndexSize"`
	IndexSizes     interface{} `bson:"indexSizes" json:"indexSizes"`
	IndexDetails   interface{} `bson:"indexDetails" json:"indexDetails"`
	Ok             int         `bson:"ok" json:"ok"`
}

type WiredTiger struct {
	CreationString string        `bson:"creationString" json:"creationString"`
	Type           string        `bson:"type" json:"type"`
	Uri            string        `bson:"uri" json:"uri"`
	LSM            *LSM          `bson:"LSM" json:"LSM"`
	BlockManager   *BlockManager `bson:"block-manager" json:"block-manager"`
	Cache          *Cache        `bson:"cache" json:"cache"`
	Session        *Session      `bson:"session" json:"session"`
	Transaction    *Transaction  `bson:"transaction" json:"transaction"`
}

type LSM struct {
	Tree         int `bson:"bloom filters in the LSM tree" json:"bloom filters in the LSM tree"`
	Positive     int `bson:"bloom filter false positives" json:"bloom filter false positives"`
	Hits         int `bson:"bloom filter hits" json:"bloom filter hits"`
	Misses       int `bson:"bloom filter misses" json:"bloom filter misses"`
	EvictedCache int `bson:"bloom filter pages evicted from cache" json:"bloom filter pages evicted from cache"`
	ReadCache    int `bson:"bloom filter pages read into cache" json:"bloom filter pages read into cache"`
	TotalSize    int `bson:"total size of bloom filters" json:"total size of bloom filters"`
	Throttle     int `bson:"sleep for LSM checkpoint throttle" json:"sleep for LSM checkpoint throttle"`
	Chunks       int `bson:"chunks in the LSM tree" json:"chunks in the LSM tree"`
	Merge        int `bson:"highest merge generation in the LSM tree" json:"highest merge generation in the LSM tree"`
	Queries      int `bson:"queries that could have benefited from a Bloom filter that did not exist" json:"queries that could have benefited from a Bloom filter that did not exist"`
	Sleep        int `bson:"sleep for LSM merge throttle" json:"sleep for LSM merge throttle"`
}

type BlockManager struct {
	FileSize       int `bson:"file allocation unit size" json:"file allocation unit size"`
	Allocated      int `bson:"blocks allocated" json:"blocks allocated"`
	CheckpointSize int `bson:"checkpoint size" json:"checkpoint size"`
	FileExt        int `bson:"allocations requiring file extension" json:"allocations requiring file extension"`
	Freed          int `bson:"blocks freed" json:"blocks freed"`
	Magic          int `bson:"file magic number" json:"file magic number"`
	Major          int `bson:"file major version number" json:"file major version number"`
	Minor          int `bson:"minor version number" json:"minor version number"`
	Reuse          int `bson:"file bytes available for reuse" json:"file bytes available for reuse"`
	Bytes          int `bson:"file size in bytes" json:"file size in bytes"`
}

type Cache struct {
	ReadCache  int `bson:"bytes read into cache" json:"bytes read into cache"`
	WriteCache int `bson:"bytes written from cache" json:"bytes written from cache"`
	CheckPoint int `bson:"checkpoint blocked page eviction" json:"checkpoint blocked page eviction"`
	Unmodified int `bson:"unmodified pages evicted" json:"unmodified pages evicted"`

	PageDeepened    int `bson:"page split during eviction deepened the tree" json:"page split during eviction deepened the tree"`
	ModifiedEvicted int `bson:"modified pages evicted" json:"modified pages evicted"`
	DataEvicted     int `bson:"data source pages selected for eviction unable to be evicted" json:"data source pages selected for eviction unable to be evicted"`
	Eviction        int `bson:"hazard pointer blocked page eviction" json:"hazard pointer blocked page eviction"`
	PageEvicted     int `bson:"internal pages evicted" json:"internal pages evicted"`
	DurEviction     int `bson:"pages split during eviction" json:"pages split during eviction"`
	InMemory        int `bson:"in-memory page splits" json:"in-memory page splits"`
	Overflow        int `bson:"overflow values cached in memory" json:"overflow values cached in memory"`
	PageInto        int `bson:"pages read into cache" json:"pages read into cache"`
	OverflowInto    int `bson:"overflow pages read into cache" json:"overflow pages read into cache"`
	PageFrom        int `bson:"pages written from cache" json:"pages written from cache"`
}

type Session struct {
	Compact   int `bson:"object compaction" json:"object compaction"`
	CursorCnt int `bson:"open cursor count" json:"open cursor count"`
}

type Transaction struct {
	UpConflict int `bson:"update conflicts" json:"update conflicts"`
}

type CollStats struct {
	Options *options.ToolOptions

	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewCollStats(opts *options.ToolOptions, sp *db.SessionProvider) *CollStats {
	return &CollStats{
		Options:         opts,
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/collStats/#dbcmd.collStats
func (cs *CollStats) Run() (*CollectionStat, error) {
	session, err := cs.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	dest := &CollectionStat{}
	err = session.DB(cs.Options.DB).Run(bson.D{{"collStats", cs.Options.Collection}, {"scale", 1024}}, dest)

	return dest, err
}
