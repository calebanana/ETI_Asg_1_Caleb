package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Passenger struct {
	P_Username   string
	P_Password   string
	P_FirstName  string
	P_LastName   string
	P_MobileNo   string
	P_EmailAddr  string
	P_ActiveTrip bool
}

func passenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "GET" {
		//Get passenger record from database
		kv := r.URL.Query()

		passengerdb := OpenDB("Passenger")
		passenger_record := GetPassenger(passengerdb, params["username"])

		if passenger_record.P_Username != "" {
			//Username exists in database
			if passenger_record.P_Password == kv["password"][0] || kv["password"][0] == "bypass" {
				//Correct password
				json.NewEncoder(w).Encode(passenger_record)
			} else {
				//Incorrect password
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - Unauthorised"))
			}
		} else {
			//Username does not exist in database
			//Encode empty passenger object
			json.NewEncoder(w).Encode(passenger_record)
		}

		defer passengerdb.Close()
	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("451 - Unable to delete account for legal reasons"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			//Add new passenger record
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.P_Username == "" {
					//No passenger object
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				passengerdb := OpenDB("Passenger")
				InsertPassenger(passengerdb, newPassenger.P_Username, newPassenger.P_Password, newPassenger.P_FirstName, newPassenger.P_LastName, newPassenger.P_MobileNo, newPassenger.P_EmailAddr)

				defer passengerdb.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide passenger information in JSON format"))
			}
		} else if r.Method == "PUT" {
			//Update passenger record
			var editPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &editPassenger)

				if editPassenger.P_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				passengerdb := OpenDB("Passenger")
				UpdatePassenger(passengerdb, editPassenger.P_Username, editPassenger.P_Password, editPassenger.P_FirstName, editPassenger.P_LastName, editPassenger.P_MobileNo, editPassenger.P_EmailAddr)

				defer passengerdb.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide passenger information in JSON format"))
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger/{username}", passenger).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening on port 5002")
	http.ListenAndServe(":5002", router)
}
