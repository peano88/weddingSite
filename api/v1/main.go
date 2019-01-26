package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

var db DataBridge
var hb HandlerBridge
var r *mux.Router
var session *mgo.Session

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDb() error {
	var err error
	for i := 0; i < 5; i++ {
		session, err = mgo.Dial("db-api")
		if err == nil {
			return nil
		}
	}
	return err
}

func init() {
	err := connectToDb()
	logErr(err)

	db.Init(session)
	hb.Init(db)
	// check if admin is set, if not set it
	if err = secureAdmin(); err != nil {
		log.Fatal(err)
	}
	r = NewRouter()
}

func main() {

	defer session.Close()
	var wait time.Duration
	wait = 13 // to be fixed later

	// Cors enabling --> no nginx already takes care
	//handler := cors.Default().Handler(r)

	srv := &http.Server{
		Addr: "0.0.0.0:5000", // for outside container communication
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in with cors.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
