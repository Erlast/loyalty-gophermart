package models

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
)

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
