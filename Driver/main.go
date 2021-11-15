package main

import (
	"encoding/json"
	"fmt"
)

type People struct {
	Firstname string
	Lastname  string
	Details   struct {
		Height int
		Weight float32
	}
}

func main() {
	var person []People
	jsonString :=
		`[
        { 
            "firstname":"Wei-Meng",     
            "lastname":"Lee",
            "details": {
                "height":175,
                "weight":70.0
            }
        },
        { 
            "firstname":"Mickey",       
            "lastname":"Mouse",
            "details": {
                "height":105,
                "weight":85.5
            }
        }        
    ]`

	json.Unmarshal([]byte(jsonString), &person)
	for _, v := range person {
		fmt.Println(v.Firstname)
		fmt.Println(v.Lastname)
		fmt.Println(v.Details.Height)
		fmt.Println(v.Details.Weight)
	}
}
