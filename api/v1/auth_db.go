package main

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// Numeric code for each API
const (
	APIReadGuest   = 2
	APIUpdateGuest = 3
	APICreateGuest = 5
	APIReadAll     = 7
)

func checkValidityAuthCode(authCode int) bool {
	if authCode == 1 || authCode == 0 {
		return false
	}
	return (authCode%APIReadGuest) == 0 || (authCode%APICreateGuest) == 0 || (authCode%APIUpdateGuest) == 0 || (authCode%APIReadAll) == 0
}

// checkAuthCode check the API code against the
// auth code stored in the db
func checkAuthCode(APICode, authCodeStored int) bool {
	// If == 0 is trivial. If > stands, for sure the user is not authorized
	// Example: a typical user has C+U rights ==> value 6; if the API is the ReadAll then APICode is 7
	// and 7 > 6. The authorization is denied
	if APICode == 0 || APICode > authCodeStored {
		return false
	}
	return authCodeStored%APICode == 0
}

//Auth is the authorization structure. In a less primitive environment
//we should bind it to an expiration date as well
type Auth struct {
	ID       bson.ObjectId `bson:"_id"`
	Token    string        `bson:"token"`
	Valid    bool          `bson:"valid"`
	User     string        `bson:"user"`
	AuthCode int           `bson:"auth_code"`
}

//InsertAuth insert an auth object in the db
func (db *DataBridge) InsertAuth(token, user string, authCode int) error {
	if token == "" {
		return errors.New("No token provided")
	}

	if user == "" {
		return errors.New("No user provided")
	}

	if !checkValidityAuthCode(authCode) {
		return errors.New("AuthCode not valid")
	}

	a := Auth{
		ID:       bson.NewObjectId(),
		Token:    token,
		Valid:    true,
		User:     user,
		AuthCode: authCode,
	}

	if err := db.TokenColl.Insert(&a); err != nil {
		return err
	}

	return nil

}

//ReadAuth : read a token from the db collection and return the Auth
func (db *DataBridge) ReadAuth(token string) (Auth, error) {
	var a Auth
	if token == "" {
		return a, errors.New("No token provided")
	}

	err := db.TokenColl.Find(bson.M{"token": token}).One(&a)
	return a, err

}

//UpdateAuth allow a selected update of an auth
func (db *DataBridge) UpdateAuth(auth *Auth) error {
	//Check the existence of the auth to update
	stored, err := db.ReadAuth(auth.Token)
	if err != nil {
		return err
	}

	//Check that we are not changing something letal
	if stored.ID.Hex() != auth.ID.Hex() {
		return errors.New("Changing id is not possible")
	}
	if stored.User != auth.User {
		return errors.New("Changing user is not possible")
	}

	return db.TokenColl.UpdateId(auth.ID, bson.M{"$set": bson.M{
		"valid":     auth.Valid,
		"auth_code": auth.AuthCode,
	},
	})

}
