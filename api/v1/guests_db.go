package main

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
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

//UserIdentification contains the information used to identify an user
type UserIdentification struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	JwtToken string `json:"jwt_token"`
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

	gc, err := db.guestColl()

	if err != nil {
		return err
	}

	if err := gc.Insert(g); err != nil {
		return err
	}

	return nil
}

//ReadGuest is a row wrapper to fetch a guest
func (db *DataBridge) ReadGuest(m bson.M) (Guest, error) {
	var g Guest

	gc, err := db.guestColl()

	if err != nil {
		return g, err
	}

	err = gc.Find(m).One(&g)
	return g, err
}

//ReadAll is a row wrapper to fetch all guests
func (db *DataBridge) ReadAll() ([]Guest, error) {
	var gs []Guest

	gc, err := db.guestColl()

	if err != nil {
		return gs, err
	}
	err = gc.Find(bson.M{}).All(&gs)
	return gs, err
}

//UpdateGuest allows to modify certains attributes of a Guest
func (db *DataBridge) UpdateGuest(g *Guest) error {
	gc, err := db.guestColl()

	if err != nil {
		return err
	}
	// CHange only certain preselected attributes
	return gc.UpdateId(g.ID, bson.M{"$set": bson.M{
		"modification":       g.Modification,
		"confirmed":          g.Confirmed,
		"needs_accomodation": g.NeedsAccomodation,
		"food_requirements":  g.FoodRequirements,
	}})
}

func sanitizeUserName(username string) bool {
	re := regexp.MustCompile("^[a-z]+$") // username are fixed and can't be changed
	return re.MatchString(username)
}

//AuthGuest check the guest against the password
func (db *DataBridge) AuthGuest(username, password string) (*UserIdentification, error) {
	//sanitize the UserName
	if !sanitizeUserName(username) {
		return nil, errors.New("Invalid username")
	}

	//use read guest to get the candidate guest by UserName
	candidate, err := db.ReadGuest(bson.M{"user_name": username})
	if err != nil {
		return nil, err
	}
	//check that the password is the same (using bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &UserIdentification{
		UserName: candidate.UserName,
		ID:       candidate.ID.Hex(),
	}, nil
}

//No delete
