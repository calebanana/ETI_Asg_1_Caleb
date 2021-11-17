package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

func driver_web(w http.ResponseWriter, r *http.Request) {

	db := OpenDB()

	if r.Header.Get("Content-type") == "application/json" {

		if r.Method == "GET" {
			GetRecords(db)
		}

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

func passenger_web(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../Passenger_Web/passenger_web.html"))

	tmpl.Execute(w, nil)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/driver/{username}", driver_web).Methods("GET", "PUT", "POST")

	router.HandleFunc("/passenger", passenger_web)

	fmt.Println("Listening on port 5000")
	http.ListenAndServe(":5000", router)
}
