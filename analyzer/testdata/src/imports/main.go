package main

import (
	"github.com/manuelarte/structfieldinitorder/testdata/src/imports/structone"
	"github.com/manuelarte/structfieldinitorder/testdata/src/imports/structtwo"
	"time"
)

func main() {
	_ = structone.StructOne{ // want `fields for struct "structone.StructOne" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}

	_ = structtwo.StructTwo{ // want `fields for struct "structtwo.StructTwo" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}
}
