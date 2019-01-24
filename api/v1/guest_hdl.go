package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
)

//UserAuth is a utility Structure to contain basic auth information
type UserAuth struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

//HandlerBridge is the struct used to provide the http.Handler
type HandlerBridge struct {
	db  DataBridge
	rnd *renderer.Render
}

//Init : initialize the handlerBridge
func (hb *HandlerBridge) Init(d DataBridge) {
	hb.db = d
	hb.rnd = renderer.New()
}

//
func (hb *HandlerBridge) CreateHandler(handler http.HandlerFunc, apiCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate the Request
		if err := newValidator(r, &hb.db, apiCode)(r, &hb.db); err != nil {
			hb.rnd.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		handler.ServeHTTP(w, r)
	}

}

// AddGuest return the provided created guest with the created id
func (hb *HandlerBridge) AddGuest(w http.ResponseWriter, r *http.Request) {
	log.Printf("i'm in")
	var g Guest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		log.Print(err.Error())
		hb.rnd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := hb.db.CreateGuest(&g); err != nil {
		hb.rnd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	hb.rnd.JSON(w, http.StatusOK, g)

}

//ModifyGuest Updates a guest; very similar to the create method
func (hb *HandlerBridge) ModifyGuest(w http.ResponseWriter, r *http.Request) {
	log.Printf("i'm in")
	id, ok := mux.Vars(r)["id"]

	if !ok {
		hb.rnd.JSON(w, http.StatusBadRequest, "No Guest Id")
		return
	}

	var g Guest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		hb.rnd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if id != g.ID.Hex() {
		hb.rnd.JSON(w, http.StatusBadRequest, "Invalid Guest provided")
		return
	}

	if err := hb.db.UpdateGuest(&g); err != nil {
		hb.rnd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	hb.rnd.JSON(w, http.StatusOK, g)

}

//GetGuestByUsername Look up for a guest: use the username as parameter
func (hb *HandlerBridge) GetGuestByUsername(w http.ResponseWriter, r *http.Request) {
	un, ok := mux.Vars(r)["user_name"]

	if !ok {
		hb.rnd.JSON(w, http.StatusBadRequest, "No Id provided")
		return
	}

	var g Guest
	var err error
	g, err = hb.db.ReadGuest(bson.M{"user_name": un})

	if err != nil {
		hb.rnd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	// Clear password
	g.Password = ""
	hb.rnd.JSON(w, http.StatusOK, g)
}

//GetGuestAll retrieve all guests
func (hb *HandlerBridge) GetGuestAll(w http.ResponseWriter, r *http.Request) {

	gs, err := hb.db.ReadAll()

	if err != nil {
		hb.rnd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	hb.rnd.JSON(w, http.StatusOK, gs)
}
func createToken(ui *UserIdentification) error {
	//create the token as library object
	token := jwt.New(jwt.SigningMethodHS256)
	//fill it with claim
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = ui.ID
	claims["user"] = ui.UserName

	//set the string signed with the 'secret'
	secret, err := getSecret()
	if err != nil {
		return err
	}

	ui.JwtToken, err = token.SignedString([]byte(secret))
	return err
}

//AuthorizeGuest is used to give authorization to a guest
func (hb *HandlerBridge) AuthorizeGuest(w http.ResponseWriter, r *http.Request) {

	//Check user identification
	var u UserAuth
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		hb.rnd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	usID := hb.db.AuthGuest(u.UserName, u.Password)

	if usID.UserName == "" {
		hb.rnd.JSON(w, http.StatusBadRequest, "Not authorized")
		return
	}

	//Create a token and set it in the User Identification Structure
	if err := createToken(&usID); err != nil {
		hb.rnd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Store the created token in the db
	// Default permission: ReadOne and Update
	var err error
	if u.UserName == "admin" {
		err = hb.db.InsertAuth(usID.JwtToken, usID.UserName, ApiReadGuest*ApiUpdateGuest*ApiReadAll*ApiCreateGuest)
	} else {
		err = hb.db.InsertAuth(usID.JwtToken, usID.UserName, ApiReadGuest*ApiUpdateGuest)
	}

	if err != nil {
		hb.rnd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	hb.rnd.JSON(w, http.StatusOK, usID)
}
