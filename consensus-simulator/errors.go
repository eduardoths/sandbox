package main

import (
	"fmt"

	"github.com/google/uuid"
)

type NotFound struct{}

func (nf NotFound) Error() string {
	return "not found"
}

type InternalError struct{}

func (ie InternalError) Error() string {
	return "internal error"
}

type ReusedTransaction struct {
	TransactionID uuid.UUID
	Status        string
}

func (rt ReusedTransaction) Error() string {
	return fmt.Sprintf("reused transaction: %s is %s",
		rt.TransactionID.String(),
		rt.Status,
	)
}
