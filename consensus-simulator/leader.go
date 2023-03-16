package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/eduardoths/sandbox/consensus-simulator/internal/config"
	"github.com/eduardoths/sandbox/consensus-simulator/internal/utils/transaction"
	ihttp "github.com/eduardoths/sandbox/go-utils/http"
	workerpool "github.com/eduardoths/sandbox/go-utils/worker-pool"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Leader struct {
	app       *fiber.App
	iFollower Follower
	followers []ihttp.Client
}

func NewLeader() Leader {
	app := fiber.New()
	leader := Leader{
		app:       app,
		iFollower: NewFollower(),
		followers: make([]ihttp.Client, 0),
	}
	leader.prepareFollowers()
	leader.route()
	return leader
}

func (l Leader) ServeHTTP(port string) error {
	return l.app.Listen(port)
}

func (l Leader) route() {
	route := l.app.Use(transaction.SetToFiber, LogMiddleware)
	route.Route("/api/storage", func(router fiber.Router) {
		router.Get("/", l.handleGet)
		router.Post("/", l.handleSave)
	})
}

func (l *Leader) prepareFollowers() {
	followerURLs := config.Get().Servers.Followers
	for _, url := range followerURLs {
		httpClient := ihttp.NewClient(url, ihttp.WithTimeout(time.Second))
		l.followers = append(l.followers, httpClient)
	}
}

func (l Leader) handleGet(c *fiber.Ctx) error {
	key := c.Query("key", "")
	val, err := l.get(c.UserContext(), key)
	if err != nil {
		log.Printf("err while getting: %s", err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendString(val)
}

func (l Leader) handleSave(c *fiber.Ctx) error {
	var body SaveRequest
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Response{
			Error: []ErrorResponse{
				{Message: "Couldn't parse body"},
			},
		})
	}

	if err := l.save(ctx, body.Key, body.Value); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Error: []ErrorResponse{
				{Message: "failed"},
			},
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

func (l Leader) get(ctx context.Context, key string) (string, error) {
	return l.iFollower.Get(ctx, key)
}

func (l Leader) save(ctx context.Context, key string, value string) error {
	ctx = transaction.SetToCtx(ctx, uuid.New())
	if err := l.iFollower.save(ctx, key, value); err != nil {
		return err
	}
	okVotes := 1
	responses, err := l.createSaveRequests(ctx, key, value)
	if err != nil {
		return err
	}
	for _, resp := range responses {
		if resp.Error == nil {
			okVotes += 1
		}
	}

	minVotes := ((len(l.followers) + 1) / 2) + 1
	if okVotes >= minVotes {
		l.sendCommit(ctx)
		l.iFollower.commit(ctx)
	} else {
		l.sendRollback(ctx)
		l.iFollower.rollback(ctx)
		return InternalError{}
	}
	return nil
}

func (l Leader) createSaveRequests(ctx context.Context, key, value string) ([]*ihttp.Response, error) {
	requests := make([]workerpool.Job[*ihttp.Response], 0, len(l.followers))
	requestBody := WSRequest{
		SaveRequest: SaveRequest{
			Key:   key,
			Value: value,
		},
		Action: SAVE_ACTION,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	transactionID := transaction.GetFromCtx(ctx)
	for i := range l.followers {
		req := l.followers[i].POST(ctx, "internal-api/storage", bytes.NewReader(bodyBytes)).
			AddHeaders("Transaction-ID", transactionID.String())
		requests = append(requests, req)
	}
	pool := workerpool.NewPool[*ihttp.Response](len(l.followers))
	return pool.Start(ctx, requests), nil
}

func (l Leader) sendCommit(ctx context.Context) {
	l.iFollower.commit(ctx)
	requests := make([]workerpool.Job[*ihttp.Response], 0, len(l.followers))
	for i := range l.followers {
		requests = append(requests, l.followers[i].PUT(ctx, "internal-api/storage/commit", nil))
	}
	pool := workerpool.NewPool[*ihttp.Response](len(l.followers))
	pool.Start(ctx, requests)
}

func (l Leader) sendRollback(ctx context.Context) {
	l.iFollower.rollback(ctx)
	requests := make([]workerpool.Job[*ihttp.Response], 0, len(l.followers))
	for i := range l.followers {
		requests = append(requests, l.followers[i].PUT(ctx, "internal-api/storage/rollback", nil))
	}
	pool := workerpool.NewPool[*ihttp.Response](len(l.followers))
	pool.Start(ctx, requests)
}

func (l Leader) getFollowerAddresses() []string {
	return config.Get().Servers.Followers
}
