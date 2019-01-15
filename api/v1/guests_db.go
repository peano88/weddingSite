package main

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Guest is the structure containing information on a single
//personally invited guest
type Guest struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	Password          string        `json:"password" bson:"password"`
	Invitees          string        `json:"invitees" bson:"invitees"`
	UserName          string        `json:"user_name" bson:"user_name"`
	Country           string        `json:"country" bson:"country"`
	Language          string        `json:"language" bson:"language"`
	Modification      string        `json:"modification" bson:"modification"`
	Confirmed         bool          `json:"confirmed" bson:"confirmed"`
	NeedsAccomodation bool          `json:"needs_accomodation" bson:"needs_accomodation"`
	FoodRequirements  string        `json:"food_requirements" bson:"food_requirements"`
}

type UserIdentification struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
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

	if g.Password == "" {
		return errors.New("No password provided")
	}

	//bycrypt the password
	bytesPwd, err := bcrypt.GenerateFromPassword([]byte(g.Password), 14)

	if err != nil {
		return err
	}
	//set back the password
	g.Password = string(bytesPwd)

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
		"modification":       g.Modification,
		"confirmed":          g.Confirmed,
		"needs_accomodation": g.NeedsAccomodation,
		"food_requirements":  g.FoodRequirements,
	}})
}

//AuthGuest check the guest against the password
func (db *DataBridge) AuthGuest(username, password string) UserIdentification {
	//use read guest to get the candidate guest by UserName
	candidate, err := db.ReadGuest(bson.M{"user_name": username})
	if err != nil {
		return UserIdentification{}
	}
	//check that the password is the same (using bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(password))
	if err != nil {
		return UserIdentification{}
	}
	return UserIdentification{
		UserName: candidate.UserName,
		ID:       candidate.ID.Hex(),
	}
}

//No delete
