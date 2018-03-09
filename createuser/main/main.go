package main

import (
	"fmt"
	"mongo-tools/common/db"
	"mongo-tools/common/options"
	"mongo-tools/createuser"
	"os"
)

func main() {
	opts := options.New("mongouser")
	opts.Addrs = []string{"127.0.0.1:27017", "127.0.0.1:27018", "127.0.0.1:27019"}
	//opts.DB = "MongoReplTest"
	//opts.Source = "admin"
	//opts.Username = "root"
	//opts.Password = "123456789"
	opts.ReplicaSetName = "mongodbrepl"
	opts.Timeout = 2
	//opts.Direct = true
	opts.TCPKeepAliveSeconds = 2
	fmt.Println(opts)

	sessionProvider, err := db.NewSessionProvider(opts)
	if err != nil {
		os.Exit(-1)
	}

	//	u := &createuser.MongoUser{
	//		DbName:          "MongoReplTest",
	//		Username:        "MongoReplTest",
	//		Password:        "123456789",
	//		SessionProvider: sessionProvider,
	//	}
	//	fmt.Println(u.CreateNormalUser())
	u := &createuser.MongoUser{
		DbName:          "admin",
		Username:        "root",
		Password:        "123456789",
		SessionProvider: sessionProvider,
	}
	fmt.Println(u.CreateAdminUser())
	sessionProvider.Close()

	//	f, err := mongo.NewMongoFactoryWithDsn("mongodb://MongoReplTest:123456789@127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/MongoReplTest?replicaSet=mongodbrepl", 2)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//	s, err := f.Get()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = s.DB("MongoReplTest").C("MongoTest").Insert(bson.M{"username": "xkey", "pwd": "xxx"})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	f.Put(s)
}
