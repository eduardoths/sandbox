package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	TX_ID_METADATA = "COMMITTED_AT"
	LEADER_URL     = "LEADER_URL"
	APP_JSON       = "application/json"
)

type Follower struct {
	app                *fiber.App
	kvStorage          *Storage
	transactionStorage *Storage
}

func NewFollower() Follower {
	app := fiber.New()
	follower := Follower{
		app:                app,
		kvStorage:          NewStorage(),
		transactionStorage: NewStorage(),
	}
	follower.route()
	return follower
}

func (f Follower) ServeHTTP(port string) error {
	return f.app.Listen(port)
}

func (f Follower) handleGet(c *fiber.Ctx) error {
	key := c.Params("key")
	val, err := f.get(c.UserContext(), key)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendString(val)
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
	leaderURL := os.Getenv(LEADER_URL)
	body := SaveRequest{
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)
	_, err = http.Post(leaderURL, APP_JSON, reader)
	return err
}

func (f Follower) get(ctx context.Context, key string) (string, error) {
	data := f.kvStorage.Get(key)
	icommitID, ok := data.Metadata[TX_ID_METADATA]
	if !ok {
		return "", InternalError{}
	}

	commitID, ok := icommitID.(uuid.UUID)
	if !ok {
		return "", InternalError{}
	}

	txData := f.transactionStorage.Get(commitID.String())
	if txData.Message != COMMITTED_TX {
		return "", NotFound{}
	}
	return data.Message, nil
}
