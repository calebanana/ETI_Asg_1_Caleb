# API Reference for Microservices

## **Driver API**
The Driver API is hosted on Port 5001, and it accepts four types of requests as described in the table below. For each of the requests, it takes in the driver's username as a parameter.

For the ```GET``` function, the driver's password is also passed in through a query string to allow for authentication.

As the ```POST``` function creates a new driver account, the password is not required in the request.

For the ```PUT``` and ```DELETE``` functions, the password is also not required in the request as these commands can only be called after the driver is successfully logged in, i.e. successful authentication through the GET function.


| API Request                           | Description          |
|:--------------------------------------|:---------------------|
| ```GET /api/v1/driver/:username```    | Get driver record    |
| ```POST /api/v1/driver/:username```   | Add driver record    |
| ```PUT /api/v1/driver/:username```    | Update driver record |
| ```DELETE /api/v1/driver/:username``` | Delete driver record |

---
## **Passenger API**
The Passenger API is hosted on Port 5002, and it accepts four types of requests as described in the table below. Similar to the Driver API, each of the requests takes in the passenger's username as a parameter.

For the ```GET``` function, the passenger's password is also passed in through a query string to allow for authentication.

As the ```POST``` function creates a new passenger account, the password is not required in the request.

For the ```PUT``` and ```DELETE``` functions, the password is also not required in the request as these commands can only be called after the passenger is successfully logged in, i.e. successful authentication through the GET function.

| API Request                              | Description             |
|:-----------------------------------------|:------------------------|
| ```GET /api/v1/passenger/:username```    | Get passenger record    |
| ```POST /api/v1/passenger/:username```   | Add passenger record    |
| ```PUT /api/v1/passenger/:username```    | Update passenger record |
| ```DELETE /api/v1/passenger/:username``` | Delete passenger record |

---
## **Trip API**
The Trip API is hosted on Port 5003, and it only accepts three types of requests, ```GET```, ```POST``` and ```PUT```. It does not accept DELETE requests as trip records cannot be deleted. This API also does not take in any parameter.

| API Request               | Description        |
|:--------------------------|:-------------------|
| ```GET /api/v1/trip```    | Get trip record    |
| ```POST /api/v1/trip```   | Add trip record    |
| ```PUT /api/v1/trip```    | Update trip record |

Depending on the situation when it is used, the ```GET``` function can accept several query string parameters.

| API GET Request                                     | Description                        |
|:----------------------------------------------------|:-----------------------------------|
| ```GET /api/v1/trip?driver=:username```             | Get trip record for driver         |
| ```GET /api/v1/trip?passenger=:username```          | Get trip record for passenger      |
| ```GET /api/v1/trip?passenger=:username&all=true``` | Get all trip records for passenger |

From the table above, Trip records can be retrieved in three different methods.

When the ```driver``` parameter is present in the query string, the trip record for that driver would be retrieved. The same would happen for the ```passenger``` parameter. In the event when both the ```passenger``` and ```all``` parameters are present, all of the trips taken by that passenger would be returned.

---
[Back to main README](./README.md)