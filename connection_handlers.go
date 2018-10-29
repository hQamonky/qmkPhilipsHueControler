package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Bridge is a json object that contains information about a Hue Bridge
type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
	Username          string `json:"username"`
}

var bridges []Bridge
var chosenBridge Bridge

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

// TODO : put file in a "connection" folder
// Handles loading of page "step1.html"
func getStep1(w http.ResponseWriter, r *http.Request) {
	// Discover Hue Bridges
	html, err := getHTML("https://discovery.meethue.com/")
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
	}
	// Unmarshal string into structs
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

// Handles result of page "step1.html" after user is done
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
	fmt.Println("Selected bridge : ", r.Form.Get("bridge"))
	chosenBridge.ID = r.Form.Get("bridge")
	chosenBridge.InternalIPAddress = ""
	foundBridge := false
	for _, element := range bridges {
		if element.ID == chosenBridge.ID {
			chosenBridge.InternalIPAddress = element.InternalIPAddress
			foundBridge = true
			break
		}
	}
	// Check values for "chosenBridge"
	if !foundBridge {
		// TODO : Handle error correctly (display error message to user and redirect to step 1)
		fmt.Println("Bridge not found")
		return
	}

	//Redirect the user to the to the next step
	http.Redirect(w, r, "/connect/step2.html", http.StatusFound)
}

// Handles loading of page "step2.html"
func getStep2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing getStep2 function...")

	//Convert the bridge variable to json
	bridgeBytes, err := json.Marshal(chosenBridge)
	// Handle errors
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write the JSON chosen bridge to the response
	w.Write(bridgeBytes)
}

// Handles result of page "step2.html" after user is done
func postStep2(w http.ResponseWriter, r *http.Request) {
	// Send POST request to get the username randomly created by the bridge
	url := "http://" + chosenBridge.InternalIPAddress + "/api/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"devicetype":"qmkPhilipsHueControler#hQamonky"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := myClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check the status code is what we expect
	fmt.Println("response Status:", resp.Status)
	if resp.Status != "200 OK" {
		// TODO : Handle error
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body:", string(body))

	var jArray []interface{}
	json.Unmarshal(body, &jArray)
	jObj := jArray[0].(map[string]interface{})
	fmt.Println("jObj:", jObj)
	if jObj["error"] != nil {
		description := jObj["error"].(map[string]interface{})["description"].(string)
		fmt.Println("error description:", description)
		// TODO : Display error in popup
	} else if jObj["success"] != nil {
		chosenBridge.Username = jObj["success"].(map[string]interface{})["username"].(string)
		fmt.Println("username:", chosenBridge.Username)

		//Redirect the user to the to the next step
		http.Redirect(w, r, "/connect/step3.html", http.StatusFound)
	}
}

// Handles loading of page "step3.html"
func getStep3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing getStep3 function...")

	// Convert the bridge variable to json
	bridgeBytes, err := json.Marshal(chosenBridge)
	// Handle errors
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write the JSON chosen bridge to the response
	w.Write(bridgeBytes)

	// Store username in database
	fmt.Println("Storing username :", chosenBridge.Username)
	fmt.Println("Storing bridge :", chosenBridge)
	err = store.CreateBridge(&chosenBridge)
	if err != nil {
		fmt.Println(err)
	}
}

// Handles loading of page "step1.html"
func getBridges(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing getBridges function...")
	// // Discover Hue Bridges
	// html, err := getHTML("https://discovery.meethue.com/")
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Error: %v", err))
	// }
	// // Unmarshal string into structs
	// json.Unmarshal(html, &bridges)

	// //Convert the "bridges" variable to json
	// bridgeListBytes, err := json.Marshal(bridges)

	// New line
	bridges, err := store.GetBridges()
	// Handle errors
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		return
	}

	// Convert the bridge variable to json
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
