package models

import (
	"time"
)

type Order struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
