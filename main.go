package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Bridge is a json object that contains information about a Hue Bridge
type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getHTML(url string) ([]byte, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return html, err
}

// TODO : refacto handlers in a new "connection_handlers.go" file
func getStep1(w http.ResponseWriter, r *http.Request) {
	// Discover Hue Bridges
	html, err := getHTML("https://discovery.meethue.com/")
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
	}
	// Unmarshal string into structs
	var bridges []Bridge
	json.Unmarshal(html, &bridges)

	//Convert the "bridges" variable to json
	bridgeListBytes, err := json.Marshal(bridges)
	// Handle errors
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write the JSON list of bridges to the response
	w.Write(bridgeListBytes)
}

func postStep1(w http.ResponseWriter, r *http.Request) {
	// Get the value from the form
	err := r.ParseForm()
	// Handle errors
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the information about the bird from the form info
	// bridgeIP := r.Form.Get("ip_address")
	fmt.Println("Selected ip address : ", r.Form.Get("ip_address"))

	//Redirect the user to the to the next step
	http.Redirect(w, r, "/connect/step2.html", http.StatusFound)
}

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

	r.HandleFunc("/step1", getStep1).Methods("GET")
	r.HandleFunc("/step1", postStep1).Methods("POST")
	return r
}

func main() {
	// Create a new router
	r := newRouter()
	// Run the web server
	http.ListenAndServe(":8080", r)
}
