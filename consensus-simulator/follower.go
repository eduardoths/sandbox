package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Follower struct {
	app     *fiber.App
	storage *Storage
}

func NewFollower() Follower {
	app := fiber.New()
	follower := Follower{
		app:     app,
		storage: NewStorage(),
	}
	follower.route()
	return follower
}

func (f Follower) ServeHTTP(port string) error {
	return f.app.Listen(port)
}

func (f Follower) handleGet(c *fiber.Ctx) error {
	key := c.Params("key")
	return c.SendString(f.get(key))
}

func (f Follower) handleExternalSave(c *fiber.Ctx) error {
	var body SaveRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Response{
			Error: []ErrorResponse{
				{Message: "Couldn't parse body"},
			},
		})
	}

	f.ExternalSave(body.Key, body.Value)

	return c.SendStatus(http.StatusNoContent)
}

func (f Follower) ExternalSave(key, value string) error {
	leaderURL := os.Getenv("LEADER_URL")
	body := SaveRequest{
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)
	_, err = http.Post(leaderURL, "application/url", reader)
	return err
}

func (f Follower) get(key string) string {
	return f.storage.Get(key)
}
