package datastore

/*
	Flatten any data. If it fails, it is fatal
*/

import (
	"encoding/json"

	"github.com/phomer/scheduler/log"
)

func Serialize(input interface{}) []byte {

	output, err := json.Marshal(input)
	if err != nil {
		log.Fatal("Serialize", err, input)
	}

	return output
}

func Deserialize(input []byte, prototype interface{}) interface{} {
	if len(input) == 0 {
		return prototype
	}

	err := json.Unmarshal(input, prototype)
	if err != nil {
		log.Fatal("Deserialize", err, input)
	}

	return prototype
}
