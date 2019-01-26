package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	//"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

const (
	baseEndpointGuest = "/guests" // the endpoint
)

func testCreate(t *testing.T, provided *Guest) {

	gc, err := db.guestColl()

	if err != nil {
		t.Fatalf("Error while connecting to collection")
	}

	// Check that the guest has not already been created
	var g Guest
	gc.Find(bson.M{"user_name": "Telemachus"}).One(&g)

	// Create the request to create a guest
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(*provided); err != nil {
		t.Fatalf("Test create encode error : %s", err.Error())
	}
	req, _ := http.NewRequest("POST", baseEndpointGuest, buf)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert := assert.New(t)
	// Check status of the response
	assert.Equal(http.StatusOK, rr.Code, "check status")

	// Get directly from the db
	gc.Find(bson.M{"user_name": "Telemachus"}).One(&g)

	// Check that the guest was created
	require.NotEqual(t, bson.ObjectId(""), g.ID, "Id not correctly created")
	// Check that the stored and the provided are aligned
	assert.Equal(provided.Confirmed, g.Confirmed, "check confirmed")
	assert.Equal(provided.NeedsAccomodation, g.NeedsAccomodation, "check accomodation")
	assert.Equal(provided.Modification, g.Modification, "check modification")
	assert.Equal(provided.FoodRequirements, g.FoodRequirements, "check food requirements")
}

func testRead(t *testing.T, provided *Guest) {
	assert := assert.New(t)

	// Create the GET request
	req, _ := http.NewRequest("GET", baseEndpointGuest+"?user_name=Telemachus", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check status of the response
	require.Equal(t, http.StatusOK, rr.Code, "check status")

	var g Guest
	json.NewDecoder(rr.Body).Decode(&g)

	// Check equality (Not the ID since it is not known at this time)
	assert.Equal(provided.Invitees, g.Invitees, "wrong invitees")
	assert.Equal(provided.FoodRequirements, g.FoodRequirements, "wrong food requirements")
	assert.Equal(provided.Language, g.Language, "wrong language")
	assert.Equal(provided.Country, g.Country, "wrong country")
	assert.Equal(provided.Modification, g.Modification, "wrong modification")
	assert.Equal(provided.Confirmed, g.Confirmed, "wrong confirmed")
	assert.Equal(provided.NeedsAccomodation, g.NeedsAccomodation, "wrong needsAccomodation")

}

func testAuth(t *testing.T, provided *Guest) {
	// Test a bad password
	u := UserAuth{
		UserName: provided.UserName,
		Password: "wrongPassword",
	}

	// Create the PUT request
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&u); err != nil {
		t.Fatalf("Test auth encode error : %s", err.Error())
	}
	req, _ := http.NewRequest("POST", "/auth", buf)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert := assert.New(t)
	// Check status of the response
	assert.Equal(http.StatusBadRequest, rr.Code, "check status")

	//test the good password
	u.Password = provided.Password
	buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&u); err != nil {
		t.Fatalf("Test auth 1 encode error : %s", err.Error())
	}

	req, _ = http.NewRequest("POST", "/auth", buf)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check status of the response
	assert.Equal(http.StatusOK, rr.Code, "check ok status")

	var usID UserIdentification
	json.NewDecoder(rr.Body).Decode(&usID)

	assert.Equal(usID.UserName, provided.UserName, "check equality of UserName")

	//Retrieve the guest from the db to get the correct Id
	gc, err := db.guestColl()

	if err != nil {
		t.Fatalf("error while connecting to collection")
	}
	var g Guest
	gc.Find(bson.M{"user_name": provided.UserName}).One(&g)

	assert.Equal(g.ID.Hex(), usID.ID, "check equality of ID")

}

func testUpdate(t *testing.T, provided *Guest) {
	// Fetch the guest from db
	stored, err := db.ReadGuest(bson.M{"user_name": provided.UserName})

	// Check on stored guest
	require.Nil(t, err)
	require.NotEqual(t, bson.ObjectId(""), stored.ID, "stored with incorrect ID")

	//Change the guest attributes
	stored.Confirmed = false
	stored.Modification = "Mum is coming"
	stored.NeedsAccomodation = false
	stored.FoodRequirements = "Mum is vegetarian"

	// Create the PUT request
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&stored); err != nil {
		t.Fatalf("Test update encode error : %s", err.Error())
	}
	req, _ := http.NewRequest("PUT", baseEndpointGuest+"/"+stored.ID.Hex(), buf)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert := assert.New(t)
	// Check status of the response
	assert.Equal(http.StatusOK, rr.Code, "check status")

	// Get directly from the db
	gc, err := db.guestColl()

	if err != nil {
		t.Fatalf("error while connecting to collection")
	}
	var g Guest
	gc.Find(bson.M{"user_name": "Telemachus"}).One(&g)

	assert.Equal(stored, g, "Update not successfull")

}

func createSubTest(g *Guest, subtest func(t *testing.T, provided *Guest)) func(t *testing.T) {
	return func(t *testing.T) {
		subtest(t, g)
	}
}

/*TestMain : the main testing function. will test creation, reading and updating of
a guest */
func TestMain(t *testing.T) {
	// Open the file with the test Guest in JSON format and unmarshal it
	jsonFile, err := os.Open("../test/one_guest.json")
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var g Guest
	json.Unmarshal(byteValue, &g)

	require.Equal(t, "Telemachus", g.UserName, "Setup of test failed")

	// Run the Subtests
	t.Run("test_create", createSubTest(&g, testCreate))
	t.Run("test_auth", createSubTest(&g, testAuth))
	t.Run("test_read", createSubTest(&g, testRead))
	t.Run("test_update", createSubTest(&g, testUpdate))

	// Clean the test data

	gc, err := db.guestColl()

	if err != nil {
		t.Fatalf(err.Error())
	}
	gc.Remove(bson.M{"user_name": "Telemachus"})

}
