package models

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	ID         int64     `json:"id,omitempty"`
	UserID     int64     `json:"user_id,omitempty"`
}
