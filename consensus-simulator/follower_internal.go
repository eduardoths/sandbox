package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (f *Follower) route() {
	f.app.Route("/api", func(router fiber.Router) {
		router.Get("/storage", f.handleGet)
		router.Post("/storage", f.handleExternalSave)
	})
	f.app.Route("/ws", func(router fiber.Router) {
		router.Post("/storage", websocket.New(f.handleInternalSave))
	})
}

func (f Follower) handleInternalSave(c *websocket.Conn) {
	for {
		var body SaveRequest
		if err := c.ReadJSON(&body); err != nil {
			c.WriteJSON(Response{
				Error: []ErrorResponse{
					{Message: "Couldn't parse body"},
				},
			})
			continue
		}

		f.save(body.Key, body.Value)

		if err := c.WriteJSON(Response{
			Data: "ok",
		}); err != nil {
			f.save(body.Key, "")
			continue
		}
	}
}

func (f Follower) save(key, value string) {
	f.storage.Save(key, value)
}
