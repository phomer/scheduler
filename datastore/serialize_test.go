package datastore

import (
	"testing"

	"github.com/phomer/scheduler/log"
)

type TestStruct struct {
	Substruct *Substruct
}

type Thing struct {
	Text string
}

type Substruct struct {
	ArrayOfThings []*Thing
}

func TestSerialize(test *testing.T) {
	thing1 := &Thing{Text: "Thing1"}
	thing2 := &Thing{Text: "Thing2"}
	sub := &Substruct{ArrayOfThings: make([]*Thing, 2)}
	sub.ArrayOfThings[0] = thing1
	sub.ArrayOfThings[1] = thing2
	top := &TestStruct{Substruct: sub}

	log.Dump(top)

	buffer := Serialize(top)
	log.Dump(buffer)
}
