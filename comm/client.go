package comm

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gorilla/http"
	httpc "github.com/gorilla/http/client"
	//"net/http"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/log"
)

func NewClient() http.Client {
	return http.DefaultClient
}

func MakeRequest(config *accounts.ClientConfig, request *Request) *Response {

	url := config.GetUrl()
	client := NewClient()

	buffer := datastore.Serialize(request)
	reader := strings.NewReader(string(buffer))

	status, _, read_closer, err := client.Post(url, nil, reader)
	if err != nil {
		switch err.(type) {

		// TODO: Add in nicer error messages

		default:
			log.Fatal("Get", err)
		}
	}

	if status.Code != httpc.SUCCESS_OK {
		return NewResponse("Post Status Failed", nil)
	}

	return NewResponse("Success", read_closer)
}

// Loop until the Stream is finished.
func DisplayStream(response *Response) {

	defer response.Reader.Close()

	data, err := ioutil.ReadAll(response.Reader)
	if err != nil {
		fmt.Println("Buffer read err", err)
	} else {
		fmt.Println("Buffer:", string(data))
	}
}
