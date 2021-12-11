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
		//Get specific/all trip records
		query := r.URL.Query()

		//Get values (if any) from query string
		driver := query.Get("driver")
		passenger := query.Get("passenger")
		all := query.Get("all")

		tripdb := OpenDB("Trip")

		if all == "" {
			//Get specific trip
			var retrievedTrip Trip
			if driver != "" {
				//Get specific trip for driver
				retrievedTrip = GetDriverTrip(tripdb, driver)
			} else if passenger != "" {
				//Get specific trip for passenger
				retrievedTrip = GetPassengerTrip(tripdb, passenger)
			}
			json.NewEncoder(w).Encode(retrievedTrip)
		} else {
			//Get all past trips for passenger
			var tripArray []Trip
			tripArray = GetPastTrips(tripdb, passenger)
			json.NewEncoder(w).Encode(tripArray)
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			//Add new trip record
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newTrip)

				tripdb := OpenDB("Trip")
				driverdb := OpenDB("Driver")
				passengerdb := OpenDB("Passenger")

				//Search for available driver
				currentDriver = AssignDriver(driverdb)

				if currentDriver.D_Username == "" {
					//No available driver
					//Encode empty driver object
					json.NewEncoder(w).Encode(currentDriver)
				} else {
					//Available driver
					ChangeDriverAvailability(driverdb, currentDriver.D_Username, 0)
					currentDriver.D_IsAvailable = false

					InsertTrip(tripdb, newTrip.T_PickUpLocation, newTrip.T_DropOffLocation, newTrip.T_Passenger)
					ChangePassengerActiveTrip(passengerdb, newTrip.T_Passenger, 1)
					AddDriverToTrip(tripdb, currentDriver.D_Username, newTrip.T_Passenger)

					json.NewEncoder(w).Encode(currentDriver)

					defer tripdb.Close()
					defer driverdb.Close()
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide trip information in JSON format"))
			}
		} else if r.Method == "PUT" {
			//Update trip record (Start/End trip)
			query := r.URL.Query()

			driver := query.Get("driver")

			tripdb := OpenDB("Trip")
			driverdb := OpenDB("Driver")
			passengerdb := OpenDB("Passenger")

			retrievedTrip := GetDriverTrip(tripdb, driver)

			if retrievedTrip.T_StartDateTime == "" {
				//Start trip
				retrievedTrip.T_StartDateTime = time.Now().Format("2006-01-02 15:04:05")
				AddStartTimeToTrip(tripdb, retrievedTrip.T_StartDateTime, retrievedTrip.T_Passenger)
			} else if retrievedTrip.T_EndDateTime == "" {
				//End trip
				retrievedTrip.T_EndDateTime = time.Now().Format("2006-01-02 15:04:05")
				AddEndTimeToTrip(tripdb, retrievedTrip.T_EndDateTime, retrievedTrip.T_Passenger)
				ChangeDriverAvailability(driverdb, retrievedTrip.T_Driver, 1)
				ChangePassengerActiveTrip(passengerdb, retrievedTrip.T_Passenger, 0)
			}

			defer tripdb.Close()
			defer driverdb.Close()
			defer passengerdb.Close()

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
