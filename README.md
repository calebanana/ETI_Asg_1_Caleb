# ETI Assignment 1
## Done by 
Caleb Goh En Yu

---
## 1. Introduction
This project was created as part of an assignment for the Emerging Trends in IT (ETI) module in my diploma course. It involves the design of microservices and the implementation of REST APIs to develop and simulate a ride-sharing platform, named BackRides.

---
## 2. Design Consideration of Microservices
For this project, I have adopted Domain-Driven Design (DDD) for the architecture design of the microservices. This process consists of 4 main steps:

1. Analyse Business Domain
2. Define of Bounded Contexts
3. Define Entities and Domain Services
4. Identify Microservices

An in-depth rundown of the DDD process can be found [here](./design_considerations.md).

---
## 3. Architecture Diagram
The diagram below shows the microservice architecture that I have came up with and implemented. It consists of a front-end web application, which is the platform that users would interact with. The web application communicates with the Driver, Passenger and Trip microservices implemented using GO through HTTP requests made to the respective API endpoints. Detailed descriptions of the three APIs can be found [here](./microservices.md).

For each of the three microservices, they are able to read, write and update data onto their respective databases, implemented with MySQL.

![Architecture Diagram](BackRides_App\images\architecture_diagram.png)

## 4. Set-Up Instructions
1. Clone the Git repository.
    ```
    git clone https://github.com/calebanana/ETI_Asg_1_Caleb.git
    ```

#### **Set-Up and Create MySQL Databases**
2. Open and execute the ```BackRidesDB_SetUp.sql``` file. After executing the SQL script, three databases (DriverDB, PassengerDB and TripDB) will be created. Each of the databases contains a table corresponding to the database, i.e. Driver, Passenger and Trip respectively.


#### **Set-Up and Launch Microservices**
3. Open four Command Prompt instances and navigate to the project root folder.
    ```
    cd ETI_Asg_1_Caleb
    ```
#### **Launching Driver, Passenger and Trip Microservices**
4. On three of the Command Prompt instances, navigate to the Driver, Passenger and Trip sub-folders and launch the API GO file and its accompanying database functions files.

    ```
    cd ./Driver
    go run driverAPI.go databaseFuncs.go
    ```
    ```
    cd ./Passenger
    go run passengerAPI.go databaseFuncs.go
    ```
    ```
    cd ./Trip
    go run tripAPI.go databaseFuncs.go
    ```
    After launching the microservices successfully, the Driver, Passenger and Trip API microservices should be hosted on Port 5001, 5002 and 5003 respectively.

#### **Launching Front-End Web Application**
5. On the last Command Prompt instance, navigate to the BackRidesApp sub-folder and launch the main GO file.
    ```
    cd ./BackRidesApp
    go run main.go
    ```
6. After running the GO file, the web application should be launched on Port 3000.
    ```
    http://localhost:3000
    ```