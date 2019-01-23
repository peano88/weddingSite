package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

const (
	//GUEST_COLLECTION is the name of the MongoDb collection
	GUEST_COLLECTION = "GUESTS"
	//TOKEN_COLLECTION is the name of the MongoDb token coll
	TOKEN_COLLECTION = "TOKEN"
	//DB_NAME is the name of the MongoDb DB
	DB_NAME = "L_C_WED"
)

//DataBridge is the struct handling the Guest Collection
type DataBridge struct {
	GuestColl *mgo.Collection
	TokenColl *mgo.Collection
}

//Init initialize the connection with the MongoDb and instantiate the collection
func (db *DataBridge) Init(session *mgo.Session) {

	dataBase := session.DB(DB_NAME)
	if dataBase == nil {
		log.Fatal("Fatal error in instantiating the DB")
	}

	db.GuestColl = dataBase.C(GUEST_COLLECTION)

	if db.GuestColl == nil {
		log.Fatal("Fatal error in instantiating the guest collection")
	}

	db.TokenColl = dataBase.C(TOKEN_COLLECTION)

	if db.TokenColl == nil {
		log.Fatal("Fatal error in instantiating the token collection")
	}
}
