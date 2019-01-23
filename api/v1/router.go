package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Queries     [2]string
}

//Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		Name:        "GetGuestByUserName",
		Method:      "GET",
		Pattern:     "/guests",
		HandlerFunc: hb.CreateHandler(hb.GetGuestByUsername, ApiReadGuest),
		Queries:     [2]string{"user_name", "{user_name}"},
	},
	Route{
		Name:        "AddGuest",
		Method:      "POST",
		Pattern:     "/guests",
		HandlerFunc: hb.CreateHandler(hb.AddGuest, ApiCreateGuest),
	},
	Route{
		Name:        "UpdateGuest",
		Method:      "PUT",
		Pattern:     "/guests/{id}",
		HandlerFunc: hb.CreateHandler(hb.ModifyGuest, ApiUpdateGuest),
	},
	Route{
		Name:        "AuthGuest",
		Method:      "POST",
		Pattern:     "/auth",
		HandlerFunc: hb.AuthorizeGuest,
	},
}

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.Queries[0] != "" {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler).
				Queries(route.Queries[0], route.Queries[1])
		} else {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}

	}

	return router
}
