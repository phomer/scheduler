package log

import (
	"fmt"
	"os"
)

// Terminate now
func Fatal(message string, err error, args ...interface{}) {
	all := append([]interface{}{"ERROR:", message, err}, args...)

	fmt.Println(all...)

	os.Exit(-1)
}
