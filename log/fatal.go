package log

import "fmt"

// Terminate now
func Fatal(message string, err error, args ...interface{}) {
	all := append([]interface{}{"FATAL ERROR", message, err}, args...)

	fmt.Println(all)
	panic("Goodbye")
}
