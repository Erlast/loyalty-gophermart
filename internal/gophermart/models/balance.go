package models

type Balance struct {
	UserID         int64   `json:"user_id,omitempty"`
	CurrentBalance float32 `json:"current_balance"`
	TotalWithdrawn float32 `json:"total_withdrawn"`
}
