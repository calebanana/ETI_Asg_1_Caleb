package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "root:E50Misweakaf@tcp(127.0.0.1:3306)/BackRidesDB")

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		return db
	}
}

func GetPassengerTrip(db *sql.DB, passenger string) Trip {
	var trip Trip

	query := fmt.Sprintf("SELECT T_ID, IFNULL(T_StartDateTime, ''), IFNULL(T_EndDateTime, ''), T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger FROM Trip WHERE T_Passenger = '%s' ORDER BY T_ID DESC LIMIT 1", passenger)

	err := db.QueryRow(query).Scan(&trip.T_ID, &trip.T_StartDateTime, &trip.T_EndDateTime, &trip.T_PickUpLocation, &trip.T_DropOffLocation, &trip.T_Driver, &trip.T_Passenger)

	if err != nil {
		panic(err.Error())
	}

	return trip
}

func GetDriverTrip(db *sql.DB, driver string) Trip {
	var trip Trip

	query := fmt.Sprintf("SELECT T_ID, IFNULL(T_StartDateTime, ''), IFNULL(T_EndDateTime, ''), T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger FROM Trip WHERE T_Driver = '%s' ORDER BY T_ID DESC LIMIT 1", driver)

	err := db.QueryRow(query).Scan(&trip.T_ID, &trip.T_StartDateTime, &trip.T_EndDateTime, &trip.T_PickUpLocation, &trip.T_DropOffLocation, &trip.T_Driver, &trip.T_Passenger)

	if err != nil {
		panic(err.Error())
	}

	return trip
}

func GetPastTrips(db *sql.DB, passenger string) []Trip {
	var tripArray []Trip

	result, err := db.Query("SELECT * FROM Trip WHERE T_Passenger = ? AND NOT ISNULL(T_EndDateTime) ORDER BY T_ID DESC", passenger)

	for result.Next() {
		var trip Trip
		err = result.Scan(&trip.T_ID, &trip.T_StartDateTime, &trip.T_EndDateTime, &trip.T_PickUpLocation, &trip.T_DropOffLocation, &trip.T_Driver, &trip.T_Passenger)
		if err != nil {
			panic(err.Error())
		}
		tripArray = append(tripArray, trip)
	}
	return tripArray
}

func InsertTrip(db *sql.DB, pickuplocation string, dropofflocation string, passenger string) {
	query1 := fmt.Sprintf("INSERT INTO Trip (T_PickUpLocation, T_DropOffLocation, T_Passenger) VALUES ('%s', '%s', '%s')", pickuplocation, dropofflocation, passenger)

	_, err := db.Query(query1)

	if err != nil {
		panic(err.Error())
	}

	query2 := fmt.Sprintf("UPDATE Passenger SET P_ActiveTrip = TRUE WHERE P_Username = '%s'", passenger)
	_, err = db.Query(query2)

	if err != nil {
		panic(err.Error())
	}
}

func AssignDriver(db *sql.DB) Driver {
	var availDriver Driver
	err := db.QueryRow("SELECT * FROM Driver WHERE D_IsAvailable = TRUE LIMIT 1").Scan(&availDriver.D_Username, &availDriver.D_Password, &availDriver.D_FirstName, &availDriver.D_LastName, &availDriver.D_MobileNo, &availDriver.D_EmailAddr, &availDriver.D_NRIC, &availDriver.D_CarLicenseNo, &availDriver.D_IsAvailable)

	if err != nil {
		panic(err.Error())
	}
	return availDriver
}

func ChangeDriverAvailability(db *sql.DB, driver string, availability int) {
	query := fmt.Sprintf("UPDATE Driver SET D_IsAvailable = '%d' WHERE D_Username = '%s'", availability, driver)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func AddDriverToTrip(db *sql.DB, driver string, passenger string) {
	query := fmt.Sprintf("UPDATE Trip SET T_Driver = '%s' WHERE T_Passenger = '%s' ORDER BY T_ID DESC LIMIT 1", driver, passenger)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func AddStartTimeToTrip(db *sql.DB, startdatetime string, passenger string) {
	query := fmt.Sprintf("UPDATE Trip SET T_StartDateTime = '%s' WHERE T_Passenger = '%s' ORDER BY T_ID DESC LIMIT 1", startdatetime, passenger)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func AddEndTimeToTrip(db *sql.DB, enddatetime string, passenger string) {
	query := fmt.Sprintf("UPDATE Trip SET T_EndDateTime = '%s' WHERE T_Passenger = '%s' ORDER BY T_ID DESC LIMIT 1", enddatetime, passenger)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func ChangePassengerActiveTrip(db *sql.DB, passenger string) {
	query := fmt.Sprintf("UPDATE Passenger SET P_ActiveTrip = FALSE WHERE P_Username = '%s'", passenger)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}
