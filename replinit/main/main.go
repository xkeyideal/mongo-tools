package main

import (
	"fmt"
	"os"

	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"
	"github.com/xkeyideal/mongo-tools/replinit"
)

func main() {
	opts := options.New("replinit")
	opts.Addrs = []string{"10.100.173.192:10002"}
	//opts.DB = "admin"
	opts.Source = "admin"
	//opts.Username = "root"
	//opts.Password = "123456789"
	//opts.ReplicaSetName = "mongodbrepl"
	opts.Timeout = 2
	opts.Direct = true
	opts.TCPKeepAliveSeconds = 2
	fmt.Println(opts)

	sessionProvider, err := db.NewSessionProvider(opts)
	if err != nil {
		os.Exit(-1)
	}

	r := replinit.NewReplInitiate(sessionProvider, []string{"10.100.173.192:10002", "10.100.173.193:10002", "10.100.173.194:10002"}, "mongodbdeploy")

	fmt.Println(r.Run())
	sessionProvider.Close()
}
