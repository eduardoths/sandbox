package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type transactionCtxKey struct{}

const TRANSACTION_ID_HEADER = "Transaction-ID"

func SetTransactionIDToContext(c *fiber.Ctx) error {
	ctx := c.UserContext()
	headers := c.GetReqHeaders()
	txID, ok := headers[TRANSACTION_ID_HEADER]
	if !ok {
		txID = uuid.NewString()
	}

	if _, err := uuid.Parse(txID); err != nil {
		txID = uuid.NewString()
	}

	ctx = context.WithValue(ctx, transactionCtxKey{}, txID)
	c.SetUserContext(ctx)
	c.Response().Header.Add(TRANSACTION_ID_HEADER, txID)
	return c.Next()
}

func TransactionIDFromCtx(ctx context.Context) uuid.UUID {
	val := ctx.Value(transactionCtxKey{})
	return val.(uuid.UUID)
}
