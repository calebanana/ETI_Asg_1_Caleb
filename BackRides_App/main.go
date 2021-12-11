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

const driverURL string = "http://localhost:5001/api/v1/driver"
const passengerURL string = "http://localhost:5002/api/v1/passenger"
const tripURL string = "http://localhost:5003/api/v1/trip"

var currentDriver Driver
var currentPassenger Passenger

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

type Passenger struct {
	P_Username   string
	P_Password   string
	P_FirstName  string
	P_LastName   string
	P_MobileNo   string
	P_EmailAddr  string
	P_ActiveTrip bool
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

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func driverNewAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_new_account.html"))
		tmpl.Execute(w, nil)
	} else {
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

		_, err := http.Post(driverURL+"/"+new_driver_data.D_Username, "application/json", bytes.NewBuffer(driver_data_json))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			http.Redirect(w, r, "/driver_main", http.StatusFound)
			currentDriver = new_driver_data
		}
	}
}

func driverLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_login.html"))
		tmpl.Execute(w, nil)
	} else {
		driver_login_data := map[string]string{
			"Username": r.FormValue("d_login_username"),
			"Password": r.FormValue("d_login_password"),
		}

		response, err := http.Get(driverURL + "/" + driver_login_data["Username"] + "?password=" + driver_login_data["Password"])

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			var retrievedDriver Driver
			_ = json.Unmarshal(data, &retrievedDriver)

			if response.StatusCode != 401 {
				http.Redirect(w, r, "/driver_main", http.StatusFound)
				currentDriver = retrievedDriver
			} else {
				http.Redirect(w, r, "/driver_login", http.StatusFound)
			}
			response.Body.Close()
		}
	}
}

func driverEditAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_edit_account.html"))
		tmpl.Execute(w, currentDriver)
	} else {
		edit_driver_data := Driver{
			D_Username:     currentDriver.D_Username,
			D_Password:     r.FormValue("d_new_password"),
			D_FirstName:    r.FormValue("d_new_firstname"),
			D_LastName:     r.FormValue("d_new_lastname"),
			D_MobileNo:     r.FormValue("d_new_mobileno"),
			D_EmailAddr:    r.FormValue("d_new_emailaddr"),
			D_NRIC:         currentDriver.D_NRIC,
			D_CarLicenseNo: r.FormValue("d_new_carlicenseno"),
			D_IsAvailable:  true,
		}

		edit_driver_data_json, _ := json.Marshal(edit_driver_data)

		request, _ := http.NewRequest(http.MethodPut,
			driverURL+"/"+edit_driver_data.D_Username,
			bytes.NewBuffer(edit_driver_data_json))

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			response.Body.Close()
		}

		currentDriver = edit_driver_data
		http.Redirect(w, r, "/driver_main", http.StatusFound)
	}
}

func driverDeleteAccount(w http.ResponseWriter, r *http.Request) {
	request, _ := http.NewRequest(http.MethodDelete,
		driverURL+"/"+currentDriver.D_Username,
		nil)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	http.Redirect(w, r, "/driver_main", http.StatusFound)
}

func driverMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_main.html"))
		tmpl.Execute(w, currentDriver)
	}
}

func driverConfirmedTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))

		response1, err := http.Get(tripURL + "?driver=" + currentDriver.D_Username)

		var retrievedTrip Trip
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response1.Body)
			_ = json.Unmarshal(data, &retrievedTrip)
		}

		response2, _ := http.Get(passengerURL + "/" + retrievedTrip.T_Passenger + "?password=bypass")

		var retrievedPassenger Passenger
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response2.Body)
			_ = json.Unmarshal(data, &retrievedPassenger)
		}

		allData := map[string]interface{}{
			"passenger": retrievedPassenger,
			"driver":    currentDriver,
			"trip":      retrievedTrip,
		}
		tmpl.Execute(w, allData)
	}
}

func driverStartTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))

		request, _ := http.NewRequest(http.MethodPut,
			tripURL+"?driver="+currentDriver.D_Username,
			nil)

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response1, err := client.Do(request)

		var retrievedTrip Trip
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response1.Body)
			_ = json.Unmarshal(data, &retrievedTrip)
			response1.Body.Close()
		}

		response2, _ := http.Get(passengerURL + "/" + retrievedTrip.T_Passenger + "?password=bypass")

		var retrievedPassenger Passenger
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response2.Body)

			_ = json.Unmarshal(data, &retrievedPassenger)
		}

		allData := map[string]interface{}{
			"driver":    currentDriver,
			"passenger": retrievedPassenger,
			"trip":      retrievedTrip,
		}
		tmpl.Execute(w, allData)
	}
}

func driverEndTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))
		request, _ := http.NewRequest(http.MethodPut,
			tripURL+"?driver="+currentDriver.D_Username,
			nil)

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response1, err := client.Do(request)

		var retrievedTrip Trip
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response1.Body)
			_ = json.Unmarshal(data, &retrievedTrip)
			response1.Body.Close()
		}

		response2, _ := http.Get(passengerURL + "/" + retrievedTrip.T_Passenger + "?password=bypass")

		var retrievedPassenger Passenger
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response2.Body)

			_ = json.Unmarshal(data, &retrievedPassenger)
		}

		currentDriver.D_IsAvailable = true
		currentPassenger.P_ActiveTrip = false

		allData := map[string]interface{}{
			"driver":    currentDriver,
			"passenger": retrievedPassenger,
			"trip":      retrievedTrip,
		}
		tmpl.Execute(w, allData)
	}
}

func passengerNewAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_new_account.html"))
		tmpl.Execute(w, nil)
	} else {
		new_passenger_data := Passenger{
			P_Username:  currentPassenger.P_Username,
			P_Password:  r.FormValue("p_password"),
			P_FirstName: r.FormValue("p_firstname"),
			P_LastName:  r.FormValue("p_lastname"),
			P_MobileNo:  r.FormValue("p_mobileno"),
			P_EmailAddr: r.FormValue("p_emailaddr"),
		}
		passenger_data_json, _ := json.Marshal(new_passenger_data)

		response, err := http.Post(passengerURL+"/"+new_passenger_data.P_Username, "application/json", bytes.NewBuffer(passenger_data_json))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			response.Body.Close()

			http.Redirect(w, r, "/passenger_main", http.StatusFound)
			currentPassenger = new_passenger_data
		}
	}
}

func passengerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_login.html"))
		tmpl.Execute(w, nil)
	} else {
		passenger_login_data := map[string]string{
			"Username": r.FormValue("p_login_username"),
			"Password": r.FormValue("p_login_password"),
		}

		response, err := http.Get(passengerURL + "/" + passenger_login_data["Username"] + "?password=" + passenger_login_data["Password"])

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			var retrievedPassenger Passenger
			_ = json.Unmarshal(data, &retrievedPassenger)

			if response.StatusCode != 401 {
				http.Redirect(w, r, "/passenger_main", http.StatusFound)
				currentPassenger = retrievedPassenger
			} else {
				http.Redirect(w, r, "/passenger_login", http.StatusFound)
			}
			response.Body.Close()
		}
	}
}

func passengerEditAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_edit_account.html"))
		tmpl.Execute(w, currentPassenger)
	} else {
		edit_passenger_data := Passenger{
			P_Username:  currentPassenger.P_Username,
			P_Password:  r.FormValue("p_new_password"),
			P_FirstName: r.FormValue("p_new_firstname"),
			P_LastName:  r.FormValue("p_new_lastname"),
			P_MobileNo:  r.FormValue("p_new_mobileno"),
			P_EmailAddr: r.FormValue("p_new_emailaddr"),
		}

		edit_passenger_data_json, _ := json.Marshal(edit_passenger_data)

		request, _ := http.NewRequest(http.MethodPut,
			passengerURL+"/"+edit_passenger_data.P_Username,
			bytes.NewBuffer(edit_passenger_data_json))

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			response.Body.Close()
		}

		currentPassenger = edit_passenger_data
		http.Redirect(w, r, "/passenger_main", http.StatusFound)
	}
}

func passengerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	request, _ := http.NewRequest(http.MethodDelete,
		passengerURL+"/"+currentPassenger.P_Username,
		nil)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	http.Redirect(w, r, "/passenger_main", http.StatusFound)
}

func passengerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_main.html"))
		tmpl.Execute(w, currentPassenger)
	}
}

func passengerNewTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_new_trip.html"))
		tmpl.Execute(w, currentPassenger)
	} else {
		new_trip_data := Trip{
			T_PickUpLocation:  r.FormValue("t_pickuplocation"),
			T_DropOffLocation: r.FormValue("t_dropofflocation"),
			T_Passenger:       currentPassenger.P_Username,
		}

		new_trip_data_json, _ := json.Marshal(new_trip_data)

		response, err := http.Post(tripURL, "application/json", bytes.NewBuffer(new_trip_data_json))

		currentPassenger.P_ActiveTrip = true

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			response.Body.Close()

			http.Redirect(w, r, "/passenger_confirmed_trip", http.StatusFound)
		}
	}
}

func passengerConfirmedTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_confirmed_trip.html"))

		response1, _ := http.Get(tripURL + "?passenger=" + currentPassenger.P_Username)
		data1, _ := ioutil.ReadAll(response1.Body)

		var retrievedTrip Trip
		_ = json.Unmarshal(data1, &retrievedTrip)
		response1.Body.Close()

		response2, _ := http.Get(driverURL + "/" + retrievedTrip.T_Driver + "?password=bypass")
		data2, _ := ioutil.ReadAll(response2.Body)

		var retrievedDriver Driver
		_ = json.Unmarshal(data2, &retrievedDriver)
		response2.Body.Close()

		allData := map[string]interface{}{
			"driver":    retrievedDriver,
			"passenger": currentPassenger,
		}
		tmpl.Execute(w, allData)
	}
}

func passengerPastTrips(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_past_trips.html"))

		response, _ := http.Get(tripURL + "?passenger=" + currentPassenger.P_Username + "&all=true")
		data, _ := ioutil.ReadAll(response.Body)

		var tripArray []Trip
		_ = json.Unmarshal(data, &tripArray)
		response.Body.Close()

		allData := map[string]interface{}{
			"passenger": currentPassenger,
			"trips":     tripArray,
		}
		tmpl.Execute(w, allData)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)

	router.HandleFunc("/driver_new_account", driverNewAccount)
	router.HandleFunc("/driver_login", driverLogin)
	router.HandleFunc("/driver_edit_account", driverEditAccount)
	router.HandleFunc("/driver_delete_account", driverDeleteAccount)
	router.HandleFunc("/driver_main", driverMain)
	router.HandleFunc("/driver_confirmed_trip", driverConfirmedTrip)
	router.HandleFunc("/driver_start_trip", driverStartTrip)
	router.HandleFunc("/driver_end_trip", driverEndTrip)

	router.HandleFunc("/passenger_new_account", passengerNewAccount)
	router.HandleFunc("/passenger_login", passengerLogin)
	router.HandleFunc("/passenger_edit_account", passengerEditAccount)
	router.HandleFunc("/passenger_delete_account", passengerDeleteAccount)
	router.HandleFunc("/passenger_main", passengerMain)
	router.HandleFunc("/passenger_new_trip", passengerNewTrip)
	router.HandleFunc("/passenger_confirmed_trip", passengerConfirmedTrip)
	router.HandleFunc("/passenger_past_trips", passengerPastTrips)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", router)
}
