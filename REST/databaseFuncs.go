package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetRecords(db *sql.DB) {
	results, err := db.Query("SELECT * FROM backridesdb.Driver")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var driver Driver
		err = results.Scan(&driver.D_Username, &driver.D_Password, &driver.D_FirstName, &driver.D_LastName, &driver.D_MobileNo, &driver.D_EmailAddr, &driver.D_NRIC, &driver.D_CarLicenseNo)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(driver.D_Username, driver.D_Password, driver.D_FirstName, driver.D_LastName,
			driver.D_MobileNo, driver.D_EmailAddr, driver.D_NRIC, driver.D_CarLicenseNo)
	}
}

func OpenDB() *sql.DB {
	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, err := sql.Open("mysql", "root:E50Misweakaf@tcp(127.0.0.1:3306)/backridesdb")

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
