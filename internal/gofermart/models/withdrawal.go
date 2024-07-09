package models

import "time"

type Withdrawal struct {
	UserID      int64     `json:"user_id,omitempty"`
	Order       string    `json:"order"`
	Amount      float64   `json:"amount"`
	ProcessedAt time.Time `json:"processed_at"`
}

// WithdrawalRequest Структура используется для представления данных, которые клиент отправляет серверу при запросе
// на списание баллов с накопительного счета.
type WithdrawalRequest struct {
	UserID int64   `json:"-"` // Поле игнорируется при маршалинге/анмаршалинге JSON
	Order  string  `json:"order"`
	Amount float64 `json:"amount"`
}
