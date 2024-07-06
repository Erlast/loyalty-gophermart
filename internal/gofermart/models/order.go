package models

import (
	"time"
)

type Order struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
