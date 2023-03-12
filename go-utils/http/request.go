package http

import (
	"io"
	"net/http"
)

type Request struct {
	Endpoint    string
	Method      string
	Headers     map[string]string
	QueryParams map[string]string
	Body        io.Reader

	req *http.Request
}
