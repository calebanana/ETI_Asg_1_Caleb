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

func GetDriver(db *sql.DB, username string) Driver {
	var driver Driver
	err := db.QueryRow("SELECT * FROM Driver WHERE D_Username = ?", username).Scan(&driver.D_Username, &driver.D_Password, &driver.D_FirstName, &driver.D_LastName, &driver.D_MobileNo, &driver.D_EmailAddr, &driver.D_NRIC, &driver.D_CarLicenseNo, &driver.D_IsAvailable)

	if err != nil {
		panic(err.Error())
	}
	return driver
}

func InsertDriver(db *sql.DB, username string, password string, firstname string, lastname string, mobileno string, emailaddr string, nric string, carlicenseno string) {
	query := fmt.Sprintf("INSERT INTO Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', TRUE)", username, password, firstname, lastname, mobileno, emailaddr, nric, carlicenseno)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdateDriver(db *sql.DB, username string, password string, firstname string, lastname string, mobileno string, emailaddr string, nric string, carlicenseno string) {
	query := fmt.Sprintf("UPDATE Driver SET D_Password = '%s', D_FirstName = '%s', D_LastName = '%s', D_MobileNo = '%s', D_EmailAddr = '%s', D_NRIC = '%s', D_CarLicenseNo = '%s' WHERE D_Username = '%s'", password, firstname, lastname, mobileno, emailaddr, nric, carlicenseno, username)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}
