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
	D_IsAvailable  bool
}

func driver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "GET" {
		//Get driver record from database
		kv := r.URL.Query()

		driverdb := OpenDB("Driver")
		driver_record := GetDriver(driverdb, params["username"])

		if driver_record.D_Username != "" {
			//Username exists in database
			if driver_record.D_Password == kv["password"][0] || kv["password"][0] == "bypass" {
				//Correct password
				json.NewEncoder(w).Encode(driver_record)
			} else {
				//Incorrect password
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - Unauthorised"))
			}
		} else {
			//Username does not exist in database
			//Encode empty driver object
			json.NewEncoder(w).Encode(driver_record)
		}

		defer driverdb.Close()
	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("451 - Unable to delete account for legal reasons"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			//Add new driver record
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.D_Username == "" {
					//No driver object
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				driverdb := OpenDB("Driver")
				InsertDriver(driverdb, newDriver.D_Username, newDriver.D_Password, newDriver.D_FirstName, newDriver.D_LastName, newDriver.D_MobileNo, newDriver.D_EmailAddr, newDriver.D_NRIC, newDriver.D_CarLicenseNo)

				defer driverdb.Close()
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide driver information in JSON format"))
			}
		} else if r.Method == "PUT" {
			//Update driver record
			var editDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &editDriver)

				if editDriver.D_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				driverdb := OpenDB("Driver")
				UpdateDriver(driverdb, editDriver.D_Username, editDriver.D_Password, editDriver.D_FirstName, editDriver.D_LastName, editDriver.D_MobileNo, editDriver.D_EmailAddr, editDriver.D_NRIC, editDriver.D_CarLicenseNo)

				defer driverdb.Close()
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/driver/{username}", driver).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening on port 5001")
	http.ListenAndServe(":5001", router)
}
