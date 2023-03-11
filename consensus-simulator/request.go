package main

import "github.com/google/uuid"

type SaveRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WSRequest struct {
	SaveRequest
	Action Action `json:"action"`
}

type Action string

const (
	SAVE_ACTION     Action = "SAVE"
	ROLLBACK_ACTION Action = "ROLLBACK"
	COMMIT_ACTION   Action = "COMMIT"
)

type Response struct {
	Data  interface{}     `json:"data,omitempty"`
	Error []ErrorResponse `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Message       string    `json:"message"`
	TransactionID uuid.UUID `json:"transaction_id,omitempty"`
}
