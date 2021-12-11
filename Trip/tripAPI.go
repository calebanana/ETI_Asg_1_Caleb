package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

type Trip struct {
	T_ID              int
	T_StartDateTime   string
	T_EndDateTime     string
	T_PickUpLocation  string
	T_DropOffLocation string
	T_Driver          string
	T_Passenger       string
}

var currentDriver Driver

func trip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.URL.Query()

		driver := query.Get("driver")
		passenger := query.Get("passenger")
		all := query.Get("all")

		db := OpenDB()

		if all == "" {
			var retrievedTrip Trip
			if driver != "" {
				retrievedTrip = GetDriverTrip(db, driver)
			} else if passenger != "" {
				retrievedTrip = GetPassengerTrip(db, passenger)
			}
			json.NewEncoder(w).Encode(retrievedTrip)
		} else {
			var tripArray []Trip
			tripArray = GetPastTrips(db, passenger)
			json.NewEncoder(w).Encode(tripArray)
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newTrip)

				db := OpenDB()
				InsertTrip(db, newTrip.T_PickUpLocation, newTrip.T_DropOffLocation, newTrip.T_Passenger)

				currentDriver = AssignDriver(db)
				ChangeDriverAvailability(db, currentDriver.D_Username, 0)
				currentDriver.D_IsAvailable = false
				AddDriverToTrip(db, currentDriver.D_Username, newTrip.T_Passenger)

				defer db.Close()

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide trip information in JSON format"))
			}
		} else if r.Method == "PUT" {
			query := r.URL.Query()

			driver := query.Get("driver")

			db := OpenDB()

			retrievedTrip := GetDriverTrip(db, driver)

			if retrievedTrip.T_StartDateTime == "" {
				retrievedTrip.T_StartDateTime = time.Now().Format("2006-01-02 15:04:05")
				AddStartTimeToTrip(db, retrievedTrip.T_StartDateTime, retrievedTrip.T_Passenger)
			} else if retrievedTrip.T_EndDateTime == "" {
				retrievedTrip.T_EndDateTime = time.Now().Format("2006-01-02 15:04:05")
				AddEndTimeToTrip(db, retrievedTrip.T_EndDateTime, retrievedTrip.T_Passenger)
				ChangeDriverAvailability(db, retrievedTrip.T_Driver, 1)
				ChangePassengerActiveTrip(db, retrievedTrip.T_Passenger)
			}
			defer db.Close()

			json.NewEncoder(w).Encode(retrievedTrip)
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/trip", trip).Methods("GET", "PUT", "POST")

	fmt.Println("Listening on port 5003")
	http.ListenAndServe(":5003", router)
}
