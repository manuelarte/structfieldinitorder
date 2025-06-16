package main

import (
	. "imports/structone"
	. "imports/structtwo"
	"time"
)

func main() {
	_ = StructOne{ // want `fields for struct "StructOne" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}

	_ = StructTwo{ // want `fields for struct "StructTwo" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}

}
