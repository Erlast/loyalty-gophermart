package validators

import (
	"testing"
)

func TestValidateOrderNumberLuhn(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		want        bool
	}{
		{"Valid order number", "79927398713", true},
		{"Valid order number with spaces", "7992 7398 713", true},
		{"Invalid order number", "79927398714", false},
		{"Invalid characters", "7992a7398713", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOrderNumberLuhn(tt.orderNumber)
			if got != tt.want {
				t.Errorf("ValidateOrderNumberLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
