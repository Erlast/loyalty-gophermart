package models

import (
	"errors"
)

type Order struct {
	UUID    string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type OrderItem struct {
	UUID  string  `json:"order"`
	Goods []Items `json:"goods"`
}

func (o *OrderItem) Validate() error {
	if o.UUID == "" {
		return errors.New("order number is required")
	}

	if len(o.Goods) == 0 {
		return errors.New("no goods found")
	}
	for _, item := range o.Goods {
		if item.Description == "" {
			return errors.New("item description is required")
		}
		if item.Price <= 0 {
			return errors.New("item price must be positive")
		}
	}
	return nil
}
