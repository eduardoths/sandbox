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

	req    *http.Request
	client Client
}

func (r *Request) Do() *Response {
	if err := r.beforeRequest(); err != nil {
		return &Response{
			Error: err,
		}
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(r.req)
	if err != nil {
		return &Response{
			Error: err,
		}
	}

	response := &Response{
		Body:       res.Body,
		StatusCode: res.StatusCode,
		Headers:    r.Headers,

		Error: nil,
		res:   res,
		req:   r,
	}
	response.validate()
	return response
}

func (r *Request) AddHeaders(keysAndValues ...string) *Request {
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		value := keysAndValues[i+1]
		r.Headers[key] = value
	}
	return r
}

func (r *Request) AddQueryParams(keysAndValues ...string) *Request {
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		value := keysAndValues[i+1]
		r.QueryParams[key] = value
	}
	return r
}

func (r *Request) beforeRequest() error {
	httpReq, err := http.NewRequest(r.Method, r.Endpoint, r.Body)
	if err != nil {
		return err
	}
	r.req = httpReq

	for k, v := range r.Headers {
		r.req.Header.Add(k, v)
	}

	queryParams := r.req.URL.Query()
	for k, v := range r.QueryParams {
		queryParams.Add(k, v)
	}
	r.req.URL.RawQuery = queryParams.Encode()

	return nil
}
