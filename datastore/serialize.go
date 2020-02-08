package datastore

/*
	Flatten any data. If it fails, it is fatal
*/

import (
	"encoding/json"
	"fmt"
)

func Serialize(input interface{}) []byte {
	output, err := json.Marshal(input)
	if err != nil {
		fmt.Println("FATAL: Unable to Serialize Persistent data ", input)
		panic("Goodbye")
	}
	return output
}

func Deserialize(input []byte, prototype interface{}) interface{} {
	if len(input) == 0 {
		return prototype
	}

	err := json.Unmarshal(input, prototype)
	if err != nil {
		fmt.Println("FATAL: Unable to Deserialize Persistent data ", input)
		panic("Goodbye")
	}
	return prototype
}
