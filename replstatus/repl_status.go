package replstatus

import (
	"time"

	"github.com/xkeyideal/mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

type ReplStatus struct {
	Set                     string       `bson:"set" json:"set"`
	Date                    time.Time    `bson:"date" json:"date"`
	MyState                 int          `bson:"myState" json:"myState"`
	Term                    int64        `bson:"term" json:"term"`
	HeartbeatIntervalMillis int64        `bson:"heartbeatIntervalMillis" json:"heartbeatIntervalMillis"`
	Members                 []ReplMember `bson:"members" json:"members"`
	Ok                      int          `bson:"ok" json:"ok"`
}

type ReplMember struct {
	Id                int                 `bson:"_id" json:"_id"`
	Name              string              `bson:"name" json:"name"`
	Health            int                 `bson:"health" json:"health"`
	State             int                 `bson:"state" json:"state"`
	StateStr          string              `bson:"stateStr" json:"stateStr"`
	Uptime            int                 `bson:"uptime" json:"uptime"`
	Optime            *ReplMemberOpTime   `bson:"optime" json:"optime"`
	OptimeDate        time.Time           `bson:"optimeDate" json:"optimeDate"`
	InfoMessage       string              `bson:"infoMessage" json:"infoMessage"`
	ElectionTime      bson.MongoTimestamp `bson:"electionTime,omitempty" json:"electionTime,omitempty"`
	ElectionDate      time.Time           `bson:"electionDate,omitempty" json:"electionDate,omitempty"`
	LastHeartbeat     time.Time           `bson:"lastHeartbeat,omitempty" json:"lastHeartbeat,omitempty"`
	LastHeartbeatRecv time.Time           `bson:"lastHeartbeatRecv,omitempty" json:"lastHeartbeatRecv,omitempty"`
	PingMs            int64               `bson:"pingMs,omitempty" json:"pingMs,omitempty"`
	SyncingTo         string              `bson:"syncingTo,omitempty" json:"syncingTo,omitempty"`
	ConfigVersion     int                 `bson:"configVersion" json:"configVersion"`
	Self              bool                `bson:"self,omitempty" json:"self,omitempty"`
}

type ReplMemberOpTime struct {
	Ts bson.MongoTimestamp `bson:"ts" json:"ts"`
	T  int64               `bson:"t" json:"t"`
}

type ReplSetGetStatus struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewReplSetGetStatus(sp *db.SessionProvider) *ReplSetGetStatus {
	return &ReplSetGetStatus{
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/replSetGetStatus/
func (repl *ReplSetGetStatus) Run() (*ReplStatus, error) {
	session, err := repl.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	dest := &ReplStatus{}
	err = session.DB("admin").Run(bson.D{{"replSetGetStatus", 1}}, dest)

	return dest, err
}
