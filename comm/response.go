package comm

import "io"

type Response struct {
	Message string
	Reader  io.ReadCloser
}

func NewResponse(message string, reader io.ReadCloser) *Response {
	return &Response{
		Message: message,
		Reader:  reader,
	}
}
