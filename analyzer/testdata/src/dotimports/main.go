package main

import (
	. "dotimports/structone"
	. "dotimports/structtwo"
	"time"
)

type MyStruct struct {
	Hello string
	Bye   string
}

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

	_ = MyStruct{ // want `fields for struct "MyStruct" are not instantiated in order`
		Bye:   "adios",
		Hello: "hola",
	}

}
