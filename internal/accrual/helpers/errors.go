package helpers

import (
	"errors"
	"fmt"
)

var ErrConflict = errors.New("status 409 conflict")

type ConflictError struct {
	Err         error
	OrderNumber string
}

func (ce *ConflictError) Error() string {
	return fmt.Sprintf("Conflict Error. Order already exists: %s, Error: %v", ce.OrderNumber, ce.Err)
}
