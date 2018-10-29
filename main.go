package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux" // The libn/pq driver is used for postgres
	_ "github.com/lib/pq"
)

// newRouter creates the router and returns it
func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Declare the static file directory and point it to the directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/connect/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/connect/").Handler(staticFileHandler).Methods("GET")

	//
	// TODO : put the connection steps in "assets/connection/" and make a subrouter to handle them
	// // Subrouter for path under "/connect"
	// connectRouter := r.PathPrefix("/connect").Subrouter()
	// connectRouter.HandleFunc("/step1", step1Handler)
	//

	r.HandleFunc("/bridges", getBridges).Methods("GET")
	r.HandleFunc("/step1", getStep1).Methods("GET")
	r.HandleFunc("/step1", postStep1).Methods("POST")
	r.HandleFunc("/step2", getStep2).Methods("GET")
	r.HandleFunc("/step2", postStep2).Methods("POST")
	r.HandleFunc("/step3", getStep3).Methods("GET")
	return r
}

func main() {
	// Connect to the database
	connString := "host=localhost port=5432 user=postgres password=qmk dbname=qmkPhilipsHueController sslmode=disable"
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	// Initialize the store
	InitStore(&dbStore{db: db})

	// Create a new router
	r := newRouter()
	// Run the web server
	http.ListenAndServe(":8080", r)
}
