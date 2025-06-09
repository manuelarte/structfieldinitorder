package main 

import (
	"time"
)

type Person struct {
	Name string
	Surname string
	Birthdate time.Time
}

func main() {
	p := Person{
		Birthdate: time.Now(),
		Surname: "Doe",
		Name: "John",
	}
}

