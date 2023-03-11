package main

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (f *Follower) route() {
	router := f.app.Use(SetTransactionIDToContext)
	publicAPI := router.Group("/api")
	publicAPI.Get("/storage", f.handleGet)
	publicAPI.Post("/storage", f.handleExternalSave)
	publicAPI.Patch("/storage/rollback")

	internalAPI := router.Post("/internal-api")
	internalAPI.Post("/storage", f.handleInternalSave)
	internalAPI.Put("/storage/rollback", f.handleInternalSave)
	internalAPI.Post("/storage/commit", f.handleInternalSave)
}

func (f Follower) handleInternalSave(c *fiber.Ctx) error {
	var body WSRequest
	ctx := c.UserContext()
	txID := TransactionIDFromCtx(ctx)
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

func (f Follower) save(ctx context.Context, key, value string) error {
	transactionID := TransactionIDFromCtx(ctx)
	f.kvStorage.Save(StorageSaveStruct{
		Key: key,
		Value: StorageData{
			Message: value,
			Metadata: map[string]interface{}{
				TX_ID_METADATA: transactionID,
			},
		},
	})
	preexistingTransaction := f.transactionStorage.Get(transactionID.String())
	if !(preexistingTransaction.Message == WATING_TX || preexistingTransaction.Message == "") {
		return ReusedTransaction{
			TransactionID: transactionID,
			Status:        preexistingTransaction.Message,
		}
	}
	f.transactionStorage.Save(StorageSaveStruct{
		Key: transactionID.String(),
		Value: StorageData{
			Message: WATING_TX,
		},
	})
	return nil
}

func (f Follower) rollback(ctx context.Context) error {
	return f.changeTransactionStatus(ctx, ROLLBACKED_TX)
}

func (f Follower) commit(ctx context.Context) error {
	return f.changeTransactionStatus(ctx, COMMITTED_TX)
}

func (f Follower) changeTransactionStatus(ctx context.Context, status string) error {
	transactionID := TransactionIDFromCtx(ctx)
	tx := f.transactionStorage.Get(transactionID.String())
	if tx.Message == "" {
		return NotFound{}
	}

	if tx.Message != WATING_TX {
		return InternalError{}
	}

	f.transactionStorage.Save(StorageSaveStruct{
		Key: transactionID.String(),
		Value: StorageData{
			Message: status,
		},
	})
	return nil
}
