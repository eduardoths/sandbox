package http

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient http.Client

	defaultHeaders     map[string]string
	validateStatusFunc func(statusCode int) bool
}

func NewClient(baseURL string, opts ...Option) Client {
	config := newConfig()
	for _, opt := range opts {
		opt(config)
	}
	return Client{
		baseURL:            baseURL,
		defaultHeaders:     config.defaultHeaders,
		validateStatusFunc: config.validateStatusCodeFn,
		httpClient: http.Client{
			Timeout: config.timeout,
		},
	}
}

func (c Client) newRequest(ctx context.Context, method string, endpoint string, body io.Reader) *Request {
	return &Request{
		BaseURL:     c.baseURL,
		Endpoint:    endpoint,
		Method:      method,
		Headers:     c.defaultHeaders,
		QueryParams: make(map[string]string),
		Body:        body,
		client:      c,
	}
}

func (c Client) GET(ctx context.Context, endpoint string) *Request {
	return c.newRequest(ctx, http.MethodGet, endpoint, nil)
}

func (c Client) POST(ctx context.Context, endpoint string, body io.Reader) *Request {
	return c.newRequest(ctx, http.MethodPost, endpoint, body)
}

func (c Client) PATCH(ctx context.Context, endpoint string, body io.Reader) *Request {
	return c.newRequest(ctx, http.MethodPatch, endpoint, body)
}

func (c Client) PUT(ctx context.Context, endpoint string, body io.Reader) *Request {
	return c.newRequest(ctx, http.MethodPut, endpoint, body)
}

func (c Client) DELETE(ctx context.Context, endpoint string, body io.Reader) *Request {
	return c.newRequest(ctx, http.MethodDelete, endpoint, body)
}
