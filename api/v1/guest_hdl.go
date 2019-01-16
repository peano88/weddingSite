package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

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

//AuthorizeGuest is used to give authorization to a guest
func (hb *HandlerBridge) AuthorizeGuest(w http.ResponseWriter, r *http.Request) {
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

	hb.rnd.JSON(w, http.StatusOK, usID)
}
