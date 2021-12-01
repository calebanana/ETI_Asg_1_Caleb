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

func driver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "GET" {

		kv := r.URL.Query()

		db := OpenDB()
		driver_record := GetDriver(db, params["username"])

		if driver_record.D_Password == kv["password"][0] {
			fmt.Println("login success")
			json.NewEncoder(w).Encode(driver_record)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorised"))
		}
		defer db.Close()
	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("451 - Unable to delete account for legal reasons"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {

			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.D_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				db := OpenDB()
				InsertDriver(db, newDriver.D_Username, newDriver.D_Password, newDriver.D_FirstName, newDriver.D_LastName, newDriver.D_MobileNo, newDriver.D_EmailAddr, newDriver.D_NRIC, newDriver.D_CarLicenseNo)

				fmt.Println("inserted driver", params["username"])
				defer db.Close()
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide driver information in JSON format"))
			}
		} else if r.Method == "PUT" {
			var editDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &editDriver)

				if editDriver.D_Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				db := OpenDB()
				UpdateDriver(db, editDriver.D_Username, editDriver.D_Password, editDriver.D_FirstName, editDriver.D_LastName, editDriver.D_MobileNo, editDriver.D_EmailAddr, editDriver.D_NRIC, editDriver.D_CarLicenseNo)

				fmt.Println(editDriver)
				fmt.Println("updated driver", params["username"])
				defer db.Close()
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/driver/{username}", driver).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening on port 5001")
	http.ListenAndServe(":5001", router)
}
