package dbstat

import (
	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"

	"gopkg.in/mgo.v2/bson"
)

type DbStatsOutput struct {
	Db              string           `bson:"db" json:"db"`
	Collections     int              `bson:"collections" json:"collections"`
	Objects         int              `bson:"objects" json:"objects"`
	AvgObjSize      int              `bson:"avgObjSize" json:"avgObjSize"`
	DataSize        int              `bson:"dataSize" json:"dataSize"`
	StorageSize     int              `bson:"storageSize" json:"storageSize"`
	NumExtents      int              `bson:"numExtents" json:"numExtents"`
	Indexes         int              `bson:"indexes" json:"indexes"`
	IndexSize       int              `bson:"indexSize" json:"indexSize"`
	FileSize        int              `bson:"fileSize" json:"fileSize"`
	NsSizeMB        int              `bson:"nsSizeMB" json:"nsSizeMB"`
	DataFileVersion *DataFileVersion `bson:"dataFileVersion" json:"dataFileVersion"`
	ExtentFreeList  *ExtentFreeList  `bson:"extentFreeList" json:"extentFreeList"`
	Ok              int              `bson:"ok" json:"ok"`
}

type DataFileVersion struct {
	Major int `bson:"major" json:"major"`
	Minor int `bson:"minor" json:"minor"`
}
type ExtentFreeList struct {
	Num  int `bson:"num" json:"num"`
	Size int `bson:"size" json:"size"`
}

type DBStats struct {
	// Generic mongo tool options
	Options *options.ToolOptions

	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewDBStats(opts *options.ToolOptions, sp *db.SessionProvider) *DBStats {
	return &DBStats{
		Options:         opts,
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/dbStats/
func (ds *DBStats) Run() (*DbStatsOutput, error) {
	session, err := ds.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	dest := &DbStatsOutput{}
	err = session.DB(ds.Options.DB).Run(bson.D{{"dbStats", 1}, {"scale", 1024}}, dest)

	return dest, err
}
