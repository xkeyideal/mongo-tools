package replinit

import (
	"strings"

	"github.com/xkeyideal/mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

// See http://docs.mongodb.org/manual/reference/replica-configuration/
// for more details
type Member struct {
	// Id is a unique id for a member in a set.
	Id int `bson:"_id"`

	// Address holds the network address of the member,
	// in the form hostname:port.
	Address string `bson:"host"`

	// priority
	Priority int `bson:"priority"`

	// votes
	Votes int `bson:"votes"`

	// tags
	Tags map[string]string `bson:"tags"`
}

// Config is the document stored in mongodb that defines the servers in the
// replica set
type Config struct {
	Name    string   `bson:"_id"`
	Members []Member `bson:"members"`
}

type ReplInitiate struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider

	Hosts    []string
	ReplName string
}

func NewReplInitiate(sp *db.SessionProvider, hosts []string, replName string) *ReplInitiate {
	return &ReplInitiate{
		SessionProvider: sp,
		Hosts:           hosts,
		ReplName:        replName,
	}
}

//制作副本集之前，随便选择其中一台机器连接，执行rs.initiate()命令即可
//https://docs.mongodb.com/v3.2/reference/command/replSetInitiate/
//2018-2-26 应国内机票项目的要求，将副本集添加上tag，且支持超过7个mongo实例的副本集
func (r *ReplInitiate) Run() error {
	session, err := r.SessionProvider.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	members := []Member{}
	for i, host := range r.Hosts {
		ss := strings.Split(host, ":")
		member := Member{
			Id:      i,
			Address: host,
			Tags: map[string]string{
				"host": ss[0],
				"port": ss[1],
			},
			// 需要给上默认值
			Priority: 1,
			Votes:    1,
		}

		// 当副本集的个数大于等于7个时，只有前7个参与投票
		if i >= 7 {
			member.Priority = 0
			member.Votes = 0
		}

		members = append(members, member)
	}

	cfg := Config{
		Name:    r.ReplName,
		Members: members,
	}
	err = session.DB("admin").Run(bson.D{{"replSetInitiate", cfg}}, nil)
	return err
}
