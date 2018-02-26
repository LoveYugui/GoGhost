package mongodb

import (
	"gopkg.in/mgo.v2"
)

var (
	UserDBCollection = Connect2User("localhost:27017")

)

func Connect2User(dbAddress string) *mgo.Collection {
	userSession, err := mgo.Dial(dbAddress)
	if err != nil {
		return nil
	}

	userDb := userSession.DB("myLord").C("user")

	return userDb
}