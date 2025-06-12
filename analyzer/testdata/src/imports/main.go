package imports

import (
	"imports/structone"
	"imports/structtwo"
	"time"
)

func main() {
	_ = structone.StructOne{ // want `fields for struct "StructOne" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}

	_ = structtwo.StructTwo{ // want `fields for struct "StructTwo" are not instantiated in order`
		Surname:   "",
		Name:      "",
		BirthDate: time.Time{},
	}
}
