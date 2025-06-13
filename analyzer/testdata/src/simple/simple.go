package main

import (
	"time"
)

type Person struct {
	Name      string
	Surname   string
	Birthdate time.Time
}

func mainSimple() {
	_ = Person{ // want `fields for struct "Person" are not instantiated in order`
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "John",
	}

	_ = &Person{ // want `fields for struct "Person" are not instantiated in order`
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "John",
	}

	_ = &Person{
		Name:      "John",
		Surname:   "Doe",
		Birthdate: time.Now(),
	}
}
