package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

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

var currentDriver Driver

const baseURL string = "http://localhost:5000/api/driver"

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	tmpl.Execute(w, nil)
}

func driverNewAccount(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("driver_new_account.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	new_driver_data := Driver{
		D_Username:     r.FormValue("d_username"),
		D_Password:     r.FormValue("d_password"),
		D_FirstName:    r.FormValue("d_firstname"),
		D_LastName:     r.FormValue("d_lastname"),
		D_MobileNo:     r.FormValue("d_mobileno"),
		D_EmailAddr:    r.FormValue("d_emailaddr"),
		D_NRIC:         r.FormValue("d_nric"),
		D_CarLicenseNo: r.FormValue("d_carlicenseno"),
	}

	driver_data_json, _ := json.Marshal(new_driver_data)

	response, err := http.Post(baseURL+"/"+new_driver_data.D_Username, "application/json", bytes.NewBuffer(driver_data_json))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	fmt.Println(string(driver_data_json))

	tmpl.Execute(w, new_driver_data)
}

func driverLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("driver_login.html"))
		tmpl.Execute(w, nil)
	} else {
		driver_login_data := map[string]string{
			"Username": r.FormValue("d_login_username"),
			"Password": r.FormValue("d_login_password"),
		}

		fmt.Println(driver_login_data)

		url := baseURL
		username := driver_login_data["Username"]
		if username != "" {
			url = baseURL + "/" + username
		}

		response, err := http.Get(url)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			var retrieved_driver Driver
			_ = json.Unmarshal(data, &retrieved_driver)

			fmt.Println(retrieved_driver)

			if driver_login_data["Password"] == retrieved_driver.D_Password {
				currentDriver = retrieved_driver
				fmt.Println("yeet")
				http.Redirect(w, r, "/driver_main", http.StatusFound)
			}
			response.Body.Close()
		}
	}
}

func driverMain(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("driver_main.html"))
	tmpl.Execute(w, currentDriver)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)

	router.HandleFunc("/driver_login", driverLogin)
	router.HandleFunc("/driver_new_account", driverNewAccount)
	router.HandleFunc("/driver_main", driverMain)

	//router.HandleFunc("/passenger_login", passengerLogin)
	//router.HandleFunc("/passenger_new_account", passengerNewAccount)
	//router.HandleFunc("/passenger_main", passengerMain)

	http.ListenAndServe(":80", router)
}
