package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/eduardoths/sandbox/consensus-simulator/internal/utils/transaction"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	TX_ID_METADATA = "TRANSACTION_ID"
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

func (f Follower) externalSave(key, value string) error {
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

func (f Follower) Get(ctx context.Context, key string) (string, error) {
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
		fmt.Println(txData)
		return "", NotFound{}
	}
	return data.Message, nil
}

func (f Follower) save(ctx context.Context, key, value string) error {
	transactionID := transaction.GetFromCtx(ctx)
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
	transactionID := transaction.GetFromCtx(ctx)
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
