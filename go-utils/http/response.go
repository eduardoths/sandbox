package http

import (
	"fmt"
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

func (r *Response) validate() {
	if r.Error != nil {
		return
	}
	if !r.req.client.validateStatusFunc(r.StatusCode) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			bodyBytes = []byte("")
		}
		r.Error = fmt.Errorf("http: request to %s at %s failed with status %d and body %v",
			r.req.client.baseURL,
			r.req.Endpoint,
			r.StatusCode,
			string(bodyBytes),
		)
	}
}
