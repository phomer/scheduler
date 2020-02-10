package comm

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"os"

	"github.com/gorilla/http"
	//"net/http"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/log"
)

func NewClient() http.Client {
	return http.DefaultClient
}

func MakeRequest(config *accounts.ClientConfig, request *Request) string {

	url := "http://127.0.0.1:8000/" // TODO: Get this from config
	client := NewClient()

	buffer := datastore.Serialize(request)
	reader := strings.NewReader(string(buffer))

	status, args, read_closer, err := client.Post(url, nil, reader)
	if err != nil {
		switch value := err.(type) {
		default:
			log.Dump("error", err)
			log.Fatal("Get", err)
			_ = value
		}
	}

	defer read_closer.Close()

	// TODO: Do something with this data
	_ = status
	_ = args

	data, err := ioutil.ReadAll(read_closer)
	if err != nil {
		fmt.Println("Buffer read err", err)
	} else {
		fmt.Println("Buffer:", string(data))
	}

	fmt.Println("Command is Sent", status)

	return "Sometext"
}

func SimpleRequest(message string) string {
	writer := os.Stdout

	url := "http://127.0.0.1:8000/"

	status, err := http.Get(writer, url)
	if err != nil {
		switch value := err.(type) {
		case *net.OpError:
			/*
				if val.Err == net.ConnectionError {
					fmt.Println("Server is not running on host")
					exit(-1)
				}
			*/
			log.Fatal("Get", err)

		default:
			log.Dump("error", err)
			log.Fatal("Get", err)
			_ = value
		}
	}

	fmt.Println("Returned status=", status)

	return "Results"
}
