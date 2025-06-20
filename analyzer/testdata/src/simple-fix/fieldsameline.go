package main

import "time"

type PersonFieldsSameLine struct {
	Name, Surname string
	Birthdate     time.Time
}

func mainComplex() {
	_ = PersonFieldsSameLine{ // want `fields for struct "PersonFieldsSameLine" are not instantiated in order`
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "John",
	}
}
