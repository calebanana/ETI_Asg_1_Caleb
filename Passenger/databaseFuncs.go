package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "root:E50Misweakaf@tcp(127.0.0.1:3306)/BackRidesDB")

	if err != nil {
		panic(err.Error())
	} else {
		return db
	}
}

func GetPassenger(db *sql.DB, username string) Passenger {
	var passenger Passenger
	err := db.QueryRow("SELECT * FROM Passenger WHERE P_Username = ?", username).Scan(&passenger.P_Username, &passenger.P_Password, &passenger.P_FirstName, &passenger.P_LastName, &passenger.P_MobileNo, &passenger.P_EmailAddr, &passenger.P_ActiveTrip)

	if err != nil {
		panic(err.Error())
	}
	return passenger
}

func InsertPassenger(db *sql.DB, username string, password string, firstname string, lastname string, mobileno string, emailaddr string) {
	query := fmt.Sprintf("INSERT INTO Passenger VALUES ('%s', '%s', '%s', '%s', '%s', '%s', FALSE)", username, password, firstname, lastname, mobileno, emailaddr)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdatePassenger(db *sql.DB, username string, password string, firstname string, lastname string, mobileno string, emailaddr string) {
	query := fmt.Sprintf("UPDATE Passenger SET P_Password = '%s', P_FirstName = '%s', P_LastName = '%s', P_MobileNo = '%s', P_EmailAddr = '%s' WHERE P_Username = '%s'", password, firstname, lastname, mobileno, emailaddr, username)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}
