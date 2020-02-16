package comm

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/gorilla/http"
	httpc "github.com/gorilla/http/client"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/jobs"
	"github.com/phomer/scheduler/log"
)

func NewClient() http.Client {
	return http.DefaultClient
}

func TokenArray(config *accounts.ClientConfig) []string {
	return []string{config.Token.Signed}
}

func MakeRequest(config *accounts.ClientConfig, request *jobs.Request) *Response {

	url := config.GetUrl(request.Type)
	client := NewClient()

	buffer := datastore.Serialize(request)
	reader := strings.NewReader(string(buffer))

	headers := map[string][]string{
		"Authorization": []string{config.Token.Signed},
		"Name":          []string{config.Username},
	}

	status, _, read_closer, err := client.Post(url, headers, reader)
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

var StopStreaming = false

func DisplayStream(response *Response) {
	defer response.Reader.Close()

	reader := bufio.NewReader(response.Reader)
	for {
		if StopStreaming {
			return
		}
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal("Streaming", err)
		}
		fmt.Printf("%s", data)
	}
}
