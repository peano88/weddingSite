package main

import (
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
	masterSession *mgo.Session
}

func (db *DataBridge) guestColl() (*mgo.Collection, error) {
	//Ping to check the availability
	if err := db.masterSession.Ping(); err != nil {
		return nil, err
	}
	copySession := db.masterSession.Clone()
	return copySession.DB(dbName).C(guestCollection), nil
}

func (db *DataBridge) tokenColl() (*mgo.Collection, error) {
	//Ping to check the availability
	if err := db.masterSession.Ping(); err != nil {
		return nil, err
	}
	copySession := db.masterSession.Clone()
	return copySession.DB(dbName).C(tokenCollection), nil
}

//Init set the master session
func (db *DataBridge) Init(session *mgo.Session) {

	db.masterSession = session

}
