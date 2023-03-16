package transaction

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type transactionCtxKey struct{}

const TRANSACTION_ID_HEADER = "Transaction-Id"

func SetToCtx(ctx context.Context, transactionID uuid.UUID) context.Context {
	return context.WithValue(ctx, transactionCtxKey{}, transactionID)
}

func SetToFiber(c *fiber.Ctx) error {
	ctx := c.UserContext()
	headers := c.GetReqHeaders()
	txIDStr, ok := headers[TRANSACTION_ID_HEADER]
	if !ok {
		txIDStr = uuid.NewString()
	}

	transactionID, err := uuid.Parse(txIDStr)
	if err != nil {
		transactionID = uuid.New()
	}

	ctx = SetToCtx(ctx, transactionID)
	c.SetUserContext(ctx)
	c.Response().Header.Add(TRANSACTION_ID_HEADER, transactionID.String())
	return c.Next()
}

func GetFromCtx(ctx context.Context) uuid.UUID {
	val := ctx.Value(transactionCtxKey{})
	id, ok := val.(uuid.UUID)
	if !ok {
		id = uuid.New()
	}
	return id
}
