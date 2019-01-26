package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func getSecret() (string, error) {
	secret, ok := os.LookupEnv("EASYWED_SECRET")
	if !ok {
		return "", errors.New("No Secret Set")
	}

	if secret == "" {
		return "", errors.New("Secret can't be empty")
	}
	return secret, nil
}

func extractToken(r *http.Request) (*jwt.Token, string, error) {
	var tokenString string

	// Get the secret as env variable
	secret, err := getSecret()
	if err != nil {
		return nil, "", err
	}
	// Get token from the Authorization header
	// format: Authorization: Bearer
	tokensString, ok := r.Header["Authorization"]
	if ok && len(tokensString) >= 1 {
		tokenString = tokensString[0]
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// If the token is empty...
	if tokenString == "" {
		// If we get here, the required token is missing
		return nil, "", errors.New("Missing auth token")
	}
	token, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, "", err
	}
	return token, tokenString, nil
}

func validateTokenForAPI(token string, APICode int, db *DataBridge) error {
	// Get the Auth Object from the db
	auth, err := db.ReadAuth(token)
	if err != nil {
		return err
	}
	if !checkAuthCode(APICode, auth.AuthCode) {
		log.Println("1")
		log.Printf("APICode:%d authCode:%d", APICode, auth.AuthCode)
		return errors.New("User is not authorized")
	}
	return nil
}

func validateRequestAPIReadGuest(r *http.Request, db *DataBridge) error {
	// Extract token
	token, tokenString, err := extractToken(r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return errors.New("Not valid Token")
	}
	// Compare the user rights
	if err = validateTokenForAPI(tokenString, APIReadGuest, db); err != nil {
		return err
	}

	// Obtain the user_name from the request
	un, ok := mux.Vars(r)["user_name"]

	if !ok {
		return errors.New("No username provided")
	}

	// Compare the user of the request with the user of the token
	user, ok := claims["user"]
	if !ok {
		return errors.New("No user claim found in token")
	}
	if user != un {
		log.Println("2")
		return errors.New("User is not authorized")
	}
	return nil
}

func validateRequestAPIReadGuestAll(r *http.Request, db *DataBridge) error {
	// Extract token
	token, tokenString, err := extractToken(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Not valid Token")
	}
	// Compare the user rights
	if err = validateTokenForAPI(tokenString, APIReadAll, db); err != nil {
		return err
	}
	return nil
}

func validateRequestAPICreateGuest(r *http.Request, db *DataBridge) error {
	// Extract token
	token, tokenString, err := extractToken(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Not valid Token")
	}
	// Compare the user rights
	if err = validateTokenForAPI(tokenString, APICreateGuest, db); err != nil {
		return err
	}

	return nil
}

func validateRequestAPIUpdateGuest(r *http.Request, db *DataBridge) error {
	// Extract token
	token, tokenString, err := extractToken(r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return errors.New("Not valid Token")
	}
	// Compare the user rights
	if err = validateTokenForAPI(tokenString, APIUpdateGuest, db); err != nil {
		return err
	}

	// Get the id from the request URI
	id, ok := mux.Vars(r)["id"]

	if !ok {
		return errors.New("No username provided")
	}

	// Unmarshal the body request: can't use directly, otherwise we will incurr
	// in an error during the real handler
	var g Guest
	var bytesBody []byte
	bytesBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	br := bytes.NewReader(bytesBody)

	if err := json.NewDecoder(br).Decode(&g); err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bytesBody))

	// first compare the id of the body with the id from the URI
	if id != g.ID.Hex() {
		return errors.New("User not authorized")
	}

	// Compare the user of the request with the user of the token
	idToken, ok := claims["id"]
	if !ok {
		return errors.New("No id claim found in token")
	}
	if id != idToken {
		return errors.New("User is not authorized")
	}
	return nil
}

func newValidator(r *http.Request, db *DataBridge, APICode int) func(r *http.Request, db *DataBridge) error {
	switch APICode {
	case APIReadGuest:
		return validateRequestAPIReadGuest
	case APIReadAll:
		return validateRequestAPIReadGuestAll
	case APICreateGuest:
		return validateRequestAPICreateGuest
	case APIUpdateGuest:
		return validateRequestAPIUpdateGuest
	}
	return func(r *http.Request, db *DataBridge) error { return nil }
}
