package serverstatus

import (
	"mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

type ServerStatus struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewServerStatus(sp *db.SessionProvider) *ServerStatus {
	return &ServerStatus{
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/serverStatus/
func (s *ServerStatus) Run() (*ServerStatusInfo, error) {
	session, err := s.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	stat := &ServerStatusInfo{}
	err = session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"metrics", 0}}, stat)

	return stat, err
}
