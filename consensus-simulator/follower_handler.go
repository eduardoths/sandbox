package main

import (
	"net/http"

	"github.com/eduardoths/sandbox/consensus-simulator/internal/utils/transaction"
	"github.com/gofiber/fiber/v2"
)

func (f Follower) ServeHTTP(port string) error {
	return f.app.Listen(port)
}

func (f *Follower) route() {
	route := f.app.Use(transaction.SetToFiber,
		LogMiddleware,
	)
	route.Route("/api/storage", func(router fiber.Router) {
		router.Get("", f.handleGet)
		router.Post("", f.handleExternalSave)
	})

	route.Route("/internal-api/storage", func(router fiber.Router) {
		route.Post("", f.handleInternalSave)
		route.Put("rollback", f.handleRollback)
		route.Put("commit", f.handleCommit)
	})
}

func (f Follower) handleGet(c *fiber.Ctx) error {
	key := c.Params("key")
	val, err := f.Get(c.UserContext(), key)
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

	f.externalSave(body.Key, body.Value)

	return c.SendStatus(http.StatusNoContent)
}

func (f Follower) handleInternalSave(c *fiber.Ctx) error {
	var body WSRequest
	ctx := c.UserContext()
	txID := transaction.GetFromCtx(ctx)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Response{
			Error: []ErrorResponse{
				{Message: "Couldn't parse body", TransactionID: txID},
			},
		})
	}

	if err := f.save(ctx, body.Key, body.Value); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Error: []ErrorResponse{
				{Message: "failed", TransactionID: txID},
			},
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

func (f Follower) handleRollback(c *fiber.Ctx) error {
	if err := f.rollback(c.UserContext()); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusNoContent)
}

func (f Follower) handleCommit(c *fiber.Ctx) error {
	if err := f.commit(c.UserContext()); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusNoContent)
}
