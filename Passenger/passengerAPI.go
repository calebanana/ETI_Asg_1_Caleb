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
	P_Username  string
	P_Password  string
	P_FirstName string
	P_LastName  string
	P_MobileNo  string
	P_EmailAddr string
}

func passenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db := OpenDB()

	if r.Method == "GET" {

		kv := r.URL.Query()
		passenger_record := GetPassenger(db, params["username"])

		if passenger_record.P_Password == kv["password"][0] {
			fmt.Println("login success")
			json.NewEncoder(w).Encode(passenger_record)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorised"))
		}

	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new passenger
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.P_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				//InsertPassenger(db, newPassenger.P_Username, newPassenger.P_Password, newPassenger.P_FirstName, newPassenger.P_LastName, newPassenger.P_MobileNo, newPassenger.P_EmailAddr)

				fmt.Println("inserted passenger", params["username"])

				// defer the close till after the main function has finished executing
				defer db.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide passenger information in JSON format"))
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/passenger/{username}", passenger).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening on port 5002")
	http.ListenAndServe(":5002", router)
}
