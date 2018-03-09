package createuser

import (
	"errors"
	"mongo-tools/common/db"

	"gopkg.in/mgo.v2"
)

const (
	mongoAdmin = "admin"
)

type MongoUser struct {
	DbName   string
	Username string
	Password string

	SessionProvider *db.SessionProvider
}

func NewMongoUser(sp *db.SessionProvider, dbName, username, password string) *MongoUser {

	return &MongoUser{
		DbName:          dbName,
		Username:        username,
		Password:        password,
		SessionProvider: sp,
	}
}

func (mu *MongoUser) CreateNormalUser() error {
	if mu.DbName == "admin" {
		return errors.New("不能对admin表添加普通用户")
	}

	user := &mgo.User{
		Username: mu.Username,
		Password: mu.Password,
		Roles:    []mgo.Role{mgo.RoleReadWrite},
	}

	return mu.SessionProvider.AddNormalUser(mu.DbName, user)
}

func (mu *MongoUser) CreateAdminUser() error {
	if mu.DbName != "admin" {
		return errors.New("超级管理员权限只能作用于admin表")
	}
	user := &mgo.User{
		Username: mu.Username,
		Password: mu.Password,
		Roles:    []mgo.Role{mgo.RoleRoot},
	}
	return mu.SessionProvider.AddAdminUser(user)
}
