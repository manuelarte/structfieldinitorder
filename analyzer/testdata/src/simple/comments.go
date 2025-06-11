package main

import "time"

type PersonComments struct {
	Name      string
	Surname   string
	Birthdate time.Time
}

func main() {
	_ = PersonComments{ // want `fields for struct "PersonComments" are not instantiated in order`
		// the birthdate is now
		Birthdate: time.Now(),
		// the surname is important
		Surname: "Doe",
		Name:    "John",
	}
}
