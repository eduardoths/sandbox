package http

import (
	"io"
	"net/http"
)

type Response struct {
	Body       io.Reader
	StatusCode int
	Headers    map[string]string

	req Request
	err error
	res *http.Response
}
