package simpleclient

import (
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	response *http.Response
}

func NewResponse(res *http.Response) *Response {
	return &Response{
		response: res,
	}
}

func (r *Response) Status() string {
	return r.response.Status
}

func (r *Response) StatusCode() int {
	return r.response.StatusCode
}

func (r *Response) Body() io.ReadCloser {
	return r.response.Body
}

func (r *Response) ContentLength() int64 {
	return r.response.ContentLength
}

func (r *Response) URL() *url.URL {
	return r.response.Request.URL
}
