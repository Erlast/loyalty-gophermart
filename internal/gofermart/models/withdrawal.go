package models

import "time"

type Withdrawal struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       string    `json:"order"`
	Amount      float64   `json:"amount"`
	UserID      int64     `json:"user_id,omitempty"`
}

// WithdrawalRequest Структура используется для представления данных, которые клиент отправляет серверу при запросе
// на списание баллов с накопительного счета.
type WithdrawalRequest struct {
	Order  string  `json:"order"`
	Amount float64 `json:"amount"`
	UserID int64   `json:"-"` // Поле игнорируется при маршалинге/анмаршалинге JSON
}
