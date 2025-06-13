package main

import (
	"time"
)

type PersonEmbedded struct {
	int
	Name      string
	Surname   string
	Birthdate time.Time
}

func mainEmbedded() {
	_ = PersonEmbedded{ // want `fields for struct "PersonEmbedded" are not instantiated in order`
		Name:      "John",
		Surname:   "Doe",
		Birthdate: time.Now(),
		int:       1,
	}

	_ = &PersonEmbedded{ // want `fields for struct "PersonEmbedded" are not instantiated in order`
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "John",
	}

	_ = &PersonEmbedded{
		Name:      "John",
		Surname:   "Doe",
		Birthdate: time.Now(),
	}
}
