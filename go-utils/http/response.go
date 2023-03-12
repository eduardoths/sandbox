package http

import (
	"io"
	"net/http"
)

type Response struct {
	Body       io.Reader
	StatusCode int
	Headers    map[string]string

	Error error
	res   *http.Response
	req   *Request
}
