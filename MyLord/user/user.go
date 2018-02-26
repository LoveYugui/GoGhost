package user

import (
	"github.com/GoGhost/persistence/mongodb"
	log "github.com/GoGhost/log"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Source  	string
	SourceId	string
	Name 		string
	Avatar		string
}

func (user *User) UserId() string {
	return user.Source + "_" + user.SourceId
}

func NewUser(source string, sourceId string, name string, avatar string) *User {
	user := &User{
		Source:source,
		SourceId:sourceId,
		Name:name,
		Avatar:avatar,
	}

	// db insert
	err := mongodb.UserDBCollection.Insert(user)

	if err != nil {
		log.Error("error insert user : ", err)
		return nil
	}

	return LoadUser(source, sourceId)
}

func LoadUser(source string, sourceId string) *User {

	var userQ User
	err := mongodb.UserDBCollection.Find(bson.M{"source":source, "sourceid":sourceId}).One(&userQ)

	if err != nil {
		log.Error("error insert user : ", err)
		return nil
	}

	return &userQ
}


