DROP DATABASE DriverDB;
DROP DATABASE PassengerDB;
DROP DATABASE TripDB;

-- Driver Database
CREATE DATABASE DriverDB;
USE DriverDB;

CREATE TABLE Driver (
    D_Username VARCHAR(20) NOT NULL,
    D_Password VARCHAR(20) NOT NULL,
    D_FirstName VARCHAR(30) NOT NULL,
    D_LastName VARCHAR(30) NOT NULL,
    D_MobileNo CHAR(8) NOT NULL,
    D_EmailAddr VARCHAR(50) NOT NULL,
    D_NRIC CHAR(9) NOT NULL,
    D_CarLicenseNo VARCHAR(8) NOT NULL,
    D_IsAvailable BOOL NOT NULL,
    PRIMARY KEY (D_Username)
);

INSERT INTO Driver
	VALUES ('dannydasimp', 'danny123', 'Danny', 'Simp', 91111111, 'dannydasimp@gmail.com', 'S1234567A', 'SBA1234D', TRUE);
INSERT INTO Driver
	VALUES ('ohhakeews', 'ohhak123', 'Oh Hak', 'Eews', 92222222, 'tokyodriftnunu@gmail.com', 'S3456789C', 'SJT9876K', TRUE);

SELECT *
	FROM Driver;

-------------------------------------------------------------------------------------------------------------

-- Passenger Database
CREATE DATABASE PassengerDB;
USE PassengerDB;

CREATE TABLE Passenger (
    P_Username VARCHAR(20) NOT NULL,
    P_Password VARCHAR(20) NOT NULL,
    P_FirstName VARCHAR(30) NOT NULL,
    P_LastName VARCHAR(30) NOT NULL,
    P_MobileNo CHAR(8) NOT NULL,
    P_EmailAddr VARCHAR(50) NOT NULL,
    P_ActiveTrip BOOL NOT NULL,
    PRIMARY KEY (P_Username)
);

INSERT INTO Passenger
	VALUES ('kennethback', 'back123', 'Kenneth', 'Back', 93333333, 'kennethisback@gmail.com', FALSE);
INSERT INTO Passenger
	VALUES ('redchicken', 'yrus123', 'Pritheev', 'Red', 94444444, 'pritheevlovesred@gmail.com', FALSE);

SELECT *
	FROM Passenger;
    
-------------------------------------------------------------------------------------------------------------

-- Trip Database
CREATE DATABASE TripDB;
USE TripDB;

CREATE TABLE Trip (
	T_ID INT NOT NULL AUTO_INCREMENT,
    T_StartDateTime DATETIME NULL,
    T_EndDateTime DATETIME NULL,
    T_PickUpLocation CHAR(6) NOT NULL,
    T_DropOffLocation CHAR(6) NOT NULL,
    T_Driver VARCHAR(20) NULL,
    T_Passenger VARCHAR(20) NOT NULL,
    PRIMARY KEY (T_ID)
);

INSERT INTO Trip (T_StartDateTime, T_EndDateTime, T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger)
	VALUES ('2021-12-07 17:12:34', '2021-12-07 18:04:56', '975314', '246809', 'ohhakeews', 'kennethback');
INSERT INTO Trip (T_StartDateTime, T_EndDateTime, T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger)
	VALUES ('2021-12-07 09:28:43', '2021-12-07 10:00:01', '145678', '435735', 'ohhakeews', 'redchicken');
INSERT INTO Trip (T_StartDateTime, T_EndDateTime, T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger)
	VALUES ('2021-12-08 12:34:56', '2021-12-08 13:04:34', '246809', '975314', 'dannydasimp', 'kennethback');
INSERT INTO Trip (T_StartDateTime, T_EndDateTime, T_PickUpLocation, T_DropOffLocation, T_Driver, T_Passenger)
	VALUES ('2021-12-08 23:45:56', '2021-12-09 00:50:12', '343434', '099887', 'dannydasimp', 'redchicken');

SELECT *
	FROM Trip;


