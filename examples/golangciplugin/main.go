package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name      string
	Surname   string
	Birthdate time.Time
}

func main() {
	p := Person{
		Birthdate: time.Now(),
		Surname:   "Doe",
		Name:      "Joe",
	}
	fmt.Printf("%+v\n", p)
}
