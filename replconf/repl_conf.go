package replconf

import (
	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"

	"gopkg.in/mgo.v2/bson"
)

type ReplConf struct {
	Id              string           `bson:"_id" json:"_id"`
	Version         int              `bson:"version" json:"version"`
	ProtocolVersion int              `bson:"protocolVersion" json:"protocolVersion"`
	Members         []ReplMember     `bson:"members" json:"members"`
	Settings        *ReplConfSetting `bson:"settings" json:"settings"`
}

type ReplMember struct {
	Id           int                    `bson:"_id" json:"_id"`
	Host         string                 `bson:"host" json:"host"`
	ArbiterOnly  bool                   `bson:"arbiterOnly" json:"arbiterOnly"`
	BuildIndexes bool                   `bson:"buildIndexes" json:"buildIndexes"`
	Hidden       bool                   `bson:"hidden" json:"hidden"`
	Priority     int                    `bson:"priority" json:"priority"`
	Tags         map[string]interface{} `bson:"tags" json:"tags"`
	SlaveDelay   int                    `bson:"slaveDelay" json:"slaveDelay"`
	Votes        int                    `bson:"votes" json:"votes"`
}

type ReplConfSetting struct {
	ChainingAllowed         bool `bson:"chainingAllowed" json:"chainingAllowed"`
	HeartbeatIntervalMillis int  `bson:"heartbeatIntervalMillis" json:"heartbeatIntervalMillis"`
	HeartbeatTimeoutSecs    int  `bson:"heartbeatTimeoutSecs" json:"heartbeatTimeoutSecs"`
	ElectionTimeoutMillis   int  `bson:"electionTimeoutMillis" json:"electionTimeoutMillis"`
}

type ReplSetGetConfig struct {
	// Generic mongo tool options
	Options *options.ToolOptions

	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func (repl *ReplSetGetConfig) Run() (*ReplConf, error) {
	session, err := repl.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	dest := &ReplConf{}
	err = session.DB("admin").Run(bson.D{{"replSetGetConfig", 1}}, dest)

	return dest, err
}
