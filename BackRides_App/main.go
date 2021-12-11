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

//API base URLs
const driverURL string = "http://localhost:5001/api/v1/driver"
const passengerURL string = "http://localhost:5002/api/v1/passenger"
const tripURL string = "http://localhost:5003/api/v1/trip"

//Local variables for storing currently logged in driver and passenger
var currentDriver Driver
var currentPassenger Passenger

//Driver, Passenger and Trip structs
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

//Landing page
func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}
}

//Create new driver account
func driverNewAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_new_account.html"))
		tmpl.Execute(w, nil)
	} else {
		//Form values
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

		//Send POST request to driver API
		_, err := http.Post(driverURL+"/"+new_driver_data.D_Username, "application/json", bytes.NewBuffer(driver_data_json))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			http.Redirect(w, r, "/driver_main", http.StatusFound)
			currentDriver = new_driver_data
		}
	}
}

//Driver login
func driverLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_login.html"))
		tmpl.Execute(w, nil)
	} else {
		//Form values
		driver_login_data := map[string]string{
			"Username": r.FormValue("d_login_username"),
			"Password": r.FormValue("d_login_password"),
		}

		//Send GET request to driver API
		response, err := http.Get(driverURL + "/" + driver_login_data["Username"] + "?password=" + driver_login_data["Password"])

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			var retrievedDriver Driver
			_ = json.Unmarshal(data, &retrievedDriver)

			if retrievedDriver.D_Username != "" {
				//Username exists in database
				//Check for authentication
				if response.StatusCode != 401 {
					//Login successful
					http.Redirect(w, r, "/driver_main", http.StatusFound)
					currentDriver = retrievedDriver
				} else {
					//Login unsuccessful
					http.Redirect(w, r, "/driver_login", http.StatusFound)
				}
				response.Body.Close()
			} else {
				//Username does not exist in database
				http.Redirect(w, r, "/driver_login", http.StatusFound)
			}
		}
	}
}

//Edit driver account
func driverEditAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_edit_account.html"))
		tmpl.Execute(w, currentDriver)
	} else {
		//Form values
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

		//Send PUT request to driver API
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

//Delete driver account
func driverDeleteAccount(w http.ResponseWriter, r *http.Request) {
	//Send DELETE request to driver API
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

//Driver main page
func driverMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_main.html"))
		tmpl.Execute(w, currentDriver)
	}
}

//Driver confirmed trip page
func driverConfirmedTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))

		//Send GET request to trip API to get current trip of driver
		response1, err := http.Get(tripURL + "?driver=" + currentDriver.D_Username)

		var retrievedTrip Trip
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response1.Body)
			_ = json.Unmarshal(data, &retrievedTrip)
		}

		//Send GET request to passenger API to get details of passenger
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

//Driver start trip
func driverStartTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))

		//Send PUT request to trip API
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

		//Send GET request to passenger API to get details of passenger
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

//Driver end trip
func driverEndTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/driver_confirmed_trip.html"))

		//Send PUT request to trip API
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

		//Send GET request to passenger API to get details of passenger
		response2, _ := http.Get(passengerURL + "/" + retrievedTrip.T_Passenger + "?password=bypass")

		var retrievedPassenger Passenger
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response2.Body)

			_ = json.Unmarshal(data, &retrievedPassenger)
		}

		//Update status of driver and passenger
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

//Create new passenger account
func passengerNewAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_new_account.html"))
		tmpl.Execute(w, nil)
	} else {
		//Form data
		new_passenger_data := Passenger{
			P_Username:  r.FormValue("p_username"),
			P_Password:  r.FormValue("p_password"),
			P_FirstName: r.FormValue("p_firstname"),
			P_LastName:  r.FormValue("p_lastname"),
			P_MobileNo:  r.FormValue("p_mobileno"),
			P_EmailAddr: r.FormValue("p_emailaddr"),
		}
		passenger_data_json, _ := json.Marshal(new_passenger_data)

		//Send POST request to passenger API
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

//Passenger login
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

			if retrievedPassenger.P_Username != "" {
				//Username exists in database
				//Check for authentication
				if response.StatusCode != 401 {
					//Login successful
					http.Redirect(w, r, "/passenger_main", http.StatusFound)
					currentPassenger = retrievedPassenger
				} else {
					//Login unsuccessful
					http.Redirect(w, r, "/passenger_login", http.StatusFound)
				}
			} else {
				//Username does not exist in database
				http.Redirect(w, r, "/passenger_login", http.StatusFound)
			}
			response.Body.Close()
		}
	}
}

//Edit passenger account
func passengerEditAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_edit_account.html"))
		tmpl.Execute(w, currentPassenger)
	} else {
		//Form data
		edit_passenger_data := Passenger{
			P_Username:  currentPassenger.P_Username,
			P_Password:  r.FormValue("p_new_password"),
			P_FirstName: r.FormValue("p_new_firstname"),
			P_LastName:  r.FormValue("p_new_lastname"),
			P_MobileNo:  r.FormValue("p_new_mobileno"),
			P_EmailAddr: r.FormValue("p_new_emailaddr"),
		}

		edit_passenger_data_json, _ := json.Marshal(edit_passenger_data)

		//Send PUT request to passenger API
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

//Delete passenger account
func passengerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	//Send DELETE request to passenger API
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

//Passenger main page
func passengerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_main.html"))
		tmpl.Execute(w, currentPassenger)
	}
}

//Passenger start new trip
func passengerNewTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_new_trip.html"))
		tmpl.Execute(w, currentPassenger)
	} else {
		//Form values
		new_trip_data := Trip{
			T_PickUpLocation:  r.FormValue("t_pickuplocation"),
			T_DropOffLocation: r.FormValue("t_dropofflocation"),
			T_Passenger:       currentPassenger.P_Username,
		}

		new_trip_data_json, _ := json.Marshal(new_trip_data)

		//Send POST request to trip API
		response, err := http.Post(tripURL, "application/json", bytes.NewBuffer(new_trip_data_json))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		data, _ := ioutil.ReadAll(response.Body)

		//Assigned driver for trip
		var retrievedDriver Driver
		_ = json.Unmarshal(data, &retrievedDriver)
		response.Body.Close()

		if retrievedDriver.D_Username == "" {
			//No available drivers
			http.Redirect(w, r, "/passenger_no_driver_found", http.StatusFound)
		} else {
			//A driver is assigned
			currentPassenger.P_ActiveTrip = true
			http.Redirect(w, r, "/passenger_confirmed_trip", http.StatusFound)
		}
	}
}

//No available drivers
func passengerNoDriverFound(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_no_driver_found.html"))
		tmpl.Execute(w, nil)
	}
}

//Passenger confirmed trip page
func passengerConfirmedTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_confirmed_trip.html"))

		//Send GET request to trip API to get current trip of passenger
		response1, _ := http.Get(tripURL + "?passenger=" + currentPassenger.P_Username)
		data1, _ := ioutil.ReadAll(response1.Body)

		var retrievedTrip Trip
		_ = json.Unmarshal(data1, &retrievedTrip)
		response1.Body.Close()

		//Send GET request to driver API to get details of driver
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

//Passenger get past trips
func passengerPastTrips(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/passenger_past_trips.html"))

		//Send GET request to trip API with "all" tag to get all trips of passenger
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
	router.HandleFunc("/passenger_no_driver_found", passengerNoDriverFound)
	router.HandleFunc("/passenger_confirmed_trip", passengerConfirmedTrip)
	router.HandleFunc("/passenger_past_trips", passengerPastTrips)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", router)
}
