package http

import "time"

type config struct {
	validateStatusCodeFn func(statusCode int) bool
	timeout              time.Duration
	defaultHeaders       map[string]string
}

func newConfig() *config {
	c := &config{}
	c.defaults()
	return c
}

func (c *config) defaults() {
	c.validateStatusCodeFn = func(statusCode int) bool {
		return statusCode >= 200 || statusCode < 300
	}
	c.defaultHeaders = map[string]string{
		"Content-Type": "application/json",
	}
	c.timeout = time.Minute
}
