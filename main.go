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
	staticFileDirectory := http.Dir("./assets/connect/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/connect/" prefix when looking for files.
	// For example, if we type "/connect/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above (/assets/connect/).
	// If we did not strip the prefix, the file server would look for "./assets/connect/connect/index.html",
	// and yield an error
	staticFileHandler := http.StripPrefix("/connect/", http.FileServer(staticFileDirectory))
	// We use a subrouter to take care of the path under /connect
	rConnectAssets := r.PathPrefix("/connect").Subrouter()
	// "/" matches "/connect/" thanks to our subrouter. Any GET request in this path
	// will route to our file server declared above
	rConnectAssets.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	// We need another subrouter who will this time route to methodes and not files
	rConnectAPI := r.PathPrefix("/api/connect").Subrouter()
	rConnectAPI.HandleFunc("/bridges", getBridges).Methods("GET")
	rConnectAPI.HandleFunc("/remove/{username}", removeBridge).Methods("GET")
	rConnectAPI.HandleFunc("/step1", getStep1).Methods("GET")
	rConnectAPI.HandleFunc("/step1", postStep1).Methods("POST")
	rConnectAPI.HandleFunc("/step2", getStep2).Methods("GET")
	rConnectAPI.HandleFunc("/step2", postStep2).Methods("POST")
	rConnectAPI.HandleFunc("/step3", getStep3).Methods("GET")
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
