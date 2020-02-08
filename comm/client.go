package comm

import (
	"fmt"

	"os"

	"github.com/gorilla/http"
)

func Request(message string) string {
	writer := os.Stdout

	url := "http://127.0.0.1:8000/"

	status, err := http.Get(writer, url)
	if err != nil {
		fmt.Println("Returned an error ", err)
		panic("Goodbye")
	}

	fmt.Println("Returned some data", status)

	return "Hey"
}
