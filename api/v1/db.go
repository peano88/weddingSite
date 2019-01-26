package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

const (
	//GUEST_COLLECTION is the name of the MongoDb collection
	guestCollection = "GUESTS"
	//TOKEN_COLLECTION is the name of the MongoDb token coll
	tokenCollection = "TOKEN"
	//DB_NAME is the name of the MongoDb DB
	dbName = "L_C_WED"
)

//DataBridge is the struct handling the Guest Collection
type DataBridge struct {
	GuestColl *mgo.Collection
	TokenColl *mgo.Collection
}

//Init initialize the connection with the MongoDb and instantiate the collection
func (db *DataBridge) Init(session *mgo.Session) {

	dataBase := session.DB(dbName)
	if dataBase == nil {
		log.Fatal("Fatal error in instantiating the DB")
	}

	db.GuestColl = dataBase.C(guestCollection)

	if db.GuestColl == nil {
		log.Fatal("Fatal error in instantiating the guest collection")
	}

	db.TokenColl = dataBase.C(tokenCollection)

	if db.TokenColl == nil {
		log.Fatal("Fatal error in instantiating the token collection")
	}
}
