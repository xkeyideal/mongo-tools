package showdbs

import (
	"mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

type DataBaseInfo struct {
	DataBases []DataBase `bson:"databases" json:"databases"`
	TotalSize int64      `bson:"totalSize" json:"totalSize"`
	Ok        int        `bson:"ok" json:"ok"`
}

type DataBase struct {
	Name       string `bson:"name" json:"name"`
	SizeOnDisk int64  `bson:"sizeOnDisk" json:"sizeOnDisk"`
	Empty      bool   `bson:"empty" json:"empty"`
}

type ShowDbs struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewShowDbs(sp *db.SessionProvider) *ShowDbs {
	return &ShowDbs{
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/listDatabases/#dbcmd.listDatabases
func (s *ShowDbs) Run() (*DataBaseInfo, error) {
	session, err := s.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	stat := &DataBaseInfo{}
	err = session.DB("admin").Run(bson.D{{"listDatabases", 1}}, stat)

	return stat, err
}
