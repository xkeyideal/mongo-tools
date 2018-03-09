package main

import (
	"encoding/json"
	"fmt"
	"mongo-tools/common/db"
	"mongo-tools/common/options"
	"mongo-tools/hostinfo"
	"os"
)

func main() {
	opts := options.New("hostinfo")
	opts.Addrs = []string{"127.0.0.1:27019"}
	//opts.Addrs = []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"}
	opts.DB = "MongoReplTest"
	opts.Source = "admin"
	opts.Username = "root"
	opts.Password = "123456789"
	//opts.ReplicaSetName = "mongodbrepl"
	opts.Timeout = 2
	opts.Direct = true
	opts.TCPKeepAliveSeconds = 2
	fmt.Println(opts)

	sessionProvider, err := db.NewSessionProvider(opts)
	if err != nil {
		os.Exit(-1)
	}

	hostinfo := hostinfo.NewHostInfo(sessionProvider)
	info, err := hostinfo.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	b, _ := json.Marshal(info)
	fmt.Println(string(b))
	sessionProvider.Close()
}
