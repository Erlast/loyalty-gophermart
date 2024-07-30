package models

type Balance struct {
	UserID         int64   `json:"user_id,omitempty"`
	CurrentBalance float32 `json:"current"`
	TotalWithdrawn float32 `json:"withdrawn"`
}
