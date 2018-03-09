package connpoolstats

import (
	"github.com/xkeyideal/mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

type ConnPoolStats struct {
	NumClientConns  int `bson:"numClientConnections" json:"numClientConnections"`
	NumAScopedConns int `bson:"numAScopedConnections" json:"numAScopedConnections"`
	TotalInUse      int `bson:"totalInUse" json:"totalInUse"`
	TotalAvail      int `bson:"totalAvailable" json:"totalAvailable"`
	TotalCreated    int `bson:"totalCreated" json:"totalCreated"`
	//TotalRefreshing int `bson:"totalRefreshing" json:"totalRefreshing"`
	//pools just support 3.2.13+ version
	//Pools           map[string]Pool `bson:"pools" json:"pools"`
	Hosts map[string]Host `bson:"hosts" json:"hosts"`
	//ReplicaSets *ReplicaSet     `bson:"replicaSets" json:"replicaSets"`
	Ok int `bson:"ok" json:"ok"`
}

type Pool struct {
	PoolInUse      int `bson:"poolInUse" json:"poolInUse"`
	PoolAvail      int `bson:"poolAvailable" json:"poolAvailable"`
	PoolCreated    int `bson:"poolCreated" json:"poolCreated"`
	PoolRefreshing int `bson:"poolRefreshing" json:"poolRefreshing"`
}

type Host struct {
	InUse   int `bson:"inUse" json:"inUse"`
	Avail   int `bson:"available" json:"available"`
	Created int `bson:"created" json:"created"`

	//refreshing just support 3.2.13+ version
	//Refreshing int `bson:"refreshing" json:"refreshing"`
}

type ReplicaSet struct {
	CsRS CsRS `bson:"csRS" json:"csRS"`
}

type CsRS struct {
	Hosts []CsRSHost `bson:"hosts" json:"hosts"`
}

type CsRSHost struct {
	Addr           string `bson:"addr" json:"addr"`
	Ok             bool   `bson:"ok" json:"ok"`
	Ismaster       bool   `bson:"ismaster" json:"ismaster"`
	Hidden         bool   `bson:"hidden" json:"hidden"`
	Secondary      bool   `bson:"secondary" json:"secondary"`
	PingTimeMillis int64  `bson:"pingTimeMillis" json:"pingTimeMillis"`
}

type ConnPoolStatsInfo struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewConnPoolStatsInfo(sp *db.SessionProvider) *ConnPoolStatsInfo {
	return &ConnPoolStatsInfo{
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/connPoolStats/#dbcmd.connPoolStats
func (s *ConnPoolStatsInfo) Run() (*ConnPoolStats, error) {
	session, err := s.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	stat := &ConnPoolStats{}
	err = session.DB("admin").Run(bson.D{{"connPoolStats", 1}}, stat)

	//	r := make(map[string]interface{})
	//	session.DB("admin").Run(bson.D{{"connPoolStats", 1}}, r)
	//	fmt.Println(r)

	return stat, err
}
