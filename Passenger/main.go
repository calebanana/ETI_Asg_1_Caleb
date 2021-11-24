package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Passenger struct {
	P_Username  string
	P_Password  string
	P_FirstName string
	P_LastName  string
	P_MobileNo  string
	P_EmailAddr string
}

var currentPassenger Passenger

const apiURL string = "http://localhost:5000/api"

func passengerNewAccount(w http.ResponseWriter, r *http.Request) {

	var newPassenger Passenger
	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &newPassenger)

		fmt.Println("newpassenger", newPassenger)

		if newPassenger.P_Username == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply passenger information in JSON format"))
			return
		} else {
			response, err := http.Post(apiURL+"/passenger/"+newPassenger.P_Username, "application/json", bytes.NewBuffer(reqBody))

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))
				response.Body.Close()
				currentPassenger = newPassenger
				http.Get("http://localhost:3000/passenger_main")
			}
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please passenger information in JSON format"))
	}

}

func passengerLogin(w http.ResponseWriter, r *http.Request) {
	var loginPassengerData map[string]string
	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &loginPassengerData)

		fmt.Println("loginpassengerdata", loginPassengerData)

		if loginPassengerData["Username"] == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply passenger information in JSON format"))
			return
		} else {
			fmt.Println(apiURL + "/passenger/" + loginPassengerData["Username"])

			response, err := http.Get(apiURL + "/passenger/" + loginPassengerData["Username"])

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))

				var retrievedPassenger Passenger
				_ = json.Unmarshal(data, &retrievedPassenger)

				fmt.Println(retrievedPassenger)

				if loginPassengerData["Password"] == retrievedPassenger.P_Password {
					currentPassenger = retrievedPassenger
					fmt.Println("correct pw")
					json.NewEncoder(w).Encode(retrievedPassenger)
					fmt.Println("sent")
				} else {
					fmt.Println("wrong pw")
					http.Get("http://localhost:3000/passenger_login")
				}
				response.Body.Close()
			}
		}
	}

}

func passengerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(currentPassenger)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/passenger_new_account", passengerNewAccount).Methods("GET", "POST")
	router.HandleFunc("/passenger_login", passengerLogin).Methods("GET", "POST")
	router.HandleFunc("/passenger_main", passengerMain)

	fmt.Println("Listening on port 3002")
	http.ListenAndServe(":3002", router)
}
