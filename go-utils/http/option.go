package http

import "time"

type Option func(*config)

func WithTimeout(t time.Duration) Option {
	return func(c *config) {
		c.timeout = t
	}
}

func WithDefaultHeaders(headers map[string]string) Option {
	return func(c *config) {
		c.defaultHeaders = headers
	}
}

func WithValidateStatusFunc(fn func(statusCode int) bool) Option {
	return func(c *config) {
		c.validateStatusCodeFn = fn
	}
}
