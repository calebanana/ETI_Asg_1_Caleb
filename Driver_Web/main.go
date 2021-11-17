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

const baseURL string = "http://localhost:5000/api/driver"

func driver_web(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("driver_web.html"))

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

func getDriver(username string) {
	url := baseURL
	if username != "" {
		url = baseURL + "/" + username
	}
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func page_2(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("page2.html"))

	tmpl.Execute(w, nil)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/driver_web", driver_web)
	router.HandleFunc("/page2", page_2)

	http.ListenAndServe(":80", router)
}
