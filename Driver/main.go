package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Driver struct {
	D_Username     string
	D_Password     string
	D_FirstName    string
	D_LastName     string
	D_MobileNo     string
	D_EmailAddr    string
	D_NRIC         string
	D_CarLicenseNo string
}

var currentDriver Driver

const apiURL string = "http://localhost:5000/api"

func driverNewAccount(w http.ResponseWriter, r *http.Request) {

	var newDriver Driver
	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &newDriver)

		fmt.Println("newdriver", newDriver)

		if newDriver.D_Username == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply driver information in JSON format"))
			return
		} else {
			response, err := http.Post(apiURL+"/driver/"+newDriver.D_Username, "application/json", bytes.NewBuffer(reqBody))

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))
				response.Body.Close()
				currentDriver = newDriver
				http.Get("http://localhost:3000/driver_main")
			}
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please driver information in JSON format"))
	}

}

func driverLogin(w http.ResponseWriter, r *http.Request) {
	var loginDriverData map[string]string
	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &loginDriverData)

		fmt.Println("logindriverdata", loginDriverData)

		if loginDriverData["Username"] == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply driver information in JSON format"))
			return
		} else {
			fmt.Println(apiURL + "/driver/" + loginDriverData["Username"])

			response, err := http.Get(apiURL + "/driver/" + loginDriverData["Username"])

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))

				var retrievedDriver Driver
				_ = json.Unmarshal(data, &retrievedDriver)

				fmt.Println(retrievedDriver)

				if loginDriverData["Password"] == retrievedDriver.D_Password {
					currentDriver = retrievedDriver
					fmt.Println("correct pw")
					json.NewEncoder(w).Encode(retrievedDriver)
					fmt.Println("sent")
				} else {
					fmt.Println("wrong pw")
					http.Get("http://localhost:3000/driver_login")
				}
				response.Body.Close()
			}
		}
	}

}

func driverMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(currentDriver)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/driver_new_account", driverNewAccount).Methods("GET", "POST")
	router.HandleFunc("/driver_login", driverLogin).Methods("GET", "POST")
	router.HandleFunc("/driver_main", driverMain)

	fmt.Println("Listening on port 3001")
	http.ListenAndServe(":3001", router)
}
