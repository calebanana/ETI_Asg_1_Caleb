package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetRecord(db *sql.DB, username string) Driver {
	var driver Driver
	err := db.QueryRow("SELECT * FROM BackRidesDB.Driver WHERE D_Username = ?", username).Scan(&driver.D_Username, &driver.D_Password, &driver.D_FirstName, &driver.D_LastName, &driver.D_MobileNo, &driver.D_EmailAddr, &driver.D_NRIC, &driver.D_CarLicenseNo)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(driver.D_Username, driver.D_Password, driver.D_FirstName, driver.D_LastName, driver.D_MobileNo, driver.D_EmailAddr, driver.D_NRIC, driver.D_CarLicenseNo)

	return driver
}

func OpenDB() *sql.DB {
	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, err := sql.Open("mysql", "root:E50Misweakaf@tcp(127.0.0.1:3306)/BackRidesDB")

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
		return db
	}
}

func InsertDriver(db *sql.DB, username string, password string, firstname string, lastname string, mobileno string, emailaddr string, nric string, carlicenseno string) {
	query := fmt.Sprintf("INSERT INTO Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", username, password, firstname, lastname, mobileno, emailaddr, nric, carlicenseno)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}
