package main

import (
	"time"
)

type Person struct {
	Name      string
	Surname   string
	Birthdate time.Time
}

func main() {
	_ = Person{ // want `fields for struct "Person" are not instantiated in order`
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "John",
	}
}
