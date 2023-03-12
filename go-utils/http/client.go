package http

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
	baseURL string

	defaultHeaders map[string]string
}

func NewClient(baseURL string) Client {
	return Client{
		baseURL: baseURL,
	}
}

func (c Client) newRequest(ctx context.Context, method string, endpoint string, body io.Reader) *Request {
	return &Request{
		Endpoint:    endpoint,
		Method:      method,
		Headers:     make(map[string]string),
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
