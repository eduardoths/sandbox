package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LogMiddleware(c *fiber.Ctx) error {
	body := c.Body()
	uri := c.Request().URI().String()

	fmt.Printf("http: received request %s %s with body %s\n", uri, c.Method(), string(body))
	if err := c.Next(); err != nil {
		return err
	}

	statusCode := c.Response().StatusCode()
	responseBody := c.Response().Body()
	fmt.Printf("http: sent response with status %d and body %s", statusCode, responseBody)
	return nil
}
