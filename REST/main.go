package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

type Passenger struct {
	P_Username  string
	P_Password  string
	P_FirstName string
	P_LastName  string
	P_MobileNo  string
	P_EmailAddr string
}

func driver(w http.ResponseWriter, r *http.Request) {

	db := OpenDB()

	if r.Method == "GET" {
		params := mux.Vars(r)

		driver_record := GetDriver(db, params["username"])

		json.NewEncoder(w).Encode(driver_record)
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new driver
		if r.Method == "POST" {

			// read the string sent to the service
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				fmt.Println(newDriver)

				if newDriver.D_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				InsertDriver(db, newDriver.D_Username, newDriver.D_Password, newDriver.D_FirstName, newDriver.D_LastName, newDriver.D_MobileNo, newDriver.D_EmailAddr, newDriver.D_NRIC, newDriver.D_CarLicenseNo)

				// defer the close till after the main function has finished executing
				defer db.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please driver information in JSON format"))
			}
		}
	}
}

func passenger(w http.ResponseWriter, r *http.Request) {

	db := OpenDB()

	if r.Method == "GET" {
		params := mux.Vars(r)

		passenger_record := GetPassenger(db, params["username"])

		json.NewEncoder(w).Encode(passenger_record)
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new driver
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				fmt.Println(newPassenger)

				if newPassenger.P_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				InsertPassenger(db, newPassenger.P_Username, newPassenger.P_Password, newPassenger.P_FirstName, newPassenger.P_LastName, newPassenger.P_MobileNo, newPassenger.P_EmailAddr)

				// defer the close till after the main function has finished executing
				defer db.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please driver information in JSON format"))
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/driver/{username}", driver).Methods("GET", "PUT", "POST")

	router.HandleFunc("/api/passenger/{username}", passenger).Methods("GET", "PUT", "POST")

	fmt.Println("Listening on port 5000")
	http.ListenAndServe(":5000", router)
}
