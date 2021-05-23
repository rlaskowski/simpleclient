package client

import "net/http"

type Response struct {
	response *http.Response
}

func NewResponse(res *http.Response) *Response {
	return &Response{
		response: res,
	}
}
