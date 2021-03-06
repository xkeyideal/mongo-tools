package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xkeyideal/mongo-tools/collnames"
	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"
)

func main() {
	opts := options.New("collstat")
	opts.Addrs = []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"}
	opts.DB = "MongoReplTest"
	opts.Collection = "MongoTest"
	opts.Source = "admin"
	opts.Username = "root"
	opts.Password = "123456789"
	opts.ReplicaSetName = "mongodbrepl"
	opts.Timeout = 2
	opts.TCPKeepAliveSeconds = 2

	sessionProvider, err := db.NewSessionProvider(opts)
	if err != nil {
		os.Exit(-1)
	}

	repl := collnames.NewCollNames(opts, sessionProvider)
	st, err := repl.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := json.Marshal(st)
	fmt.Println(string(b))
	sessionProvider.Close()
}
