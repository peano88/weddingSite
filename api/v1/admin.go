package main

import (
	"errors"
	"os"

	"gopkg.in/mgo.v2/bson"
)

func getPwd() (string, error) {
	// The password is stored as environmental variable
	pwd, ok := os.LookupEnv("EASYWED_PWD")
	if !ok {
		return "", errors.New("No Admin PWD Set")
	}

	if pwd == "" {
		return "", errors.New("Pwd can't be empty")
	}
	return pwd, nil
}

// the admin is a special guest with just username and password
func createAdmin() error {

	pwd, err := getPwd()
	if err != nil {
		return err
	}
	g := new(Guest)
	g.UserName = "admin"
	g.Password = pwd

	//Create user
	return db.CreateGuest(g)
}

func secureAdmin() error {
	// check if the admin exists
	ad, _ := db.ReadGuest(bson.M{"user_name": "admin"})

	if ad.UserName == "admin" {
		return nil
	}
	// It does not exist: create it
	return createAdmin()
}
