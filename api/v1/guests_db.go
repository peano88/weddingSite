package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Guest is the structure containing information on a single
//personally invited guest
type Guest struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	Name              string        `json:"name" bson:"name"`
	FamilyName        string        `json:"family_name" bson:"family_name"`
	UserName          string        `json:"user_name" bson:"user_name"`
	Country           string        `json:"country" bson:"country"`
	Language          string        `json:"language" bson:"language"`
	ComingWith        string        `json:"coming_with" bson:"coming_with"`
	Confirmed         bool          `json:"confirmed" bson:"confirmed"`
	NeedsAccomodation bool          `json:"needs_accomodation" bson:"needs_accomodation"`
}

const (
	//GUEST_COLLECTION is the name of the MongoDb collection
	GUEST_COLLECTION = "GUESTS"
	//DB_NAME is the name of the MongoDb DB
	DB_NAME = "L_C_WED"
)

//DataBridge is the struct handling the Guest Collection
type DataBridge struct {
	GuestColl *mgo.Collection
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
}

//CreateGuest store the provided Guest in the collection
func (db *DataBridge) CreateGuest(g *Guest) error {
	g.ID = bson.NewObjectId()

	if err := db.GuestColl.Insert(g); err != nil {
		return err
	}

	return nil
}

//ReadGuest is a row wrapper to fetch a guest
func (db *DataBridge) ReadGuest(m bson.M) (Guest, error) {
	var g Guest

	err := db.GuestColl.Find(m).One(&g)

	return g, err
}

//UpdateGuest allows to modify certains attributes of a Guest
func (db *DataBridge) UpdateGuest(g *Guest) error {
	// CHange only certain preselected attributes
	return db.GuestColl.UpdateId(g.ID, bson.M{"$set": bson.M{
		"coming_with":        g.ComingWith,
		"confirmed":          g.Confirmed,
		"needs_accomodation": g.NeedsAccomodation,
	}})
}

//No delete
