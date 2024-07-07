package models

type Balance struct {
	UserID         int64   `json:"user_id"`
	CurrentBalance float64 `json:"current_balance"`
	TotalWithdrawn float64 `json:"total_withdrawn"`
}
