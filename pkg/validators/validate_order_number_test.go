package validators

import (
	"testing"
)

func TestValidateOrderNumberLuna(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		want        bool
	}{
		{"Valid order number", "123456789012345", true},
		{"Valid order number with spaces", "1234 5678 9012 345", true},
		{"Invalid order number", "12345678901234", false},
		{"Invalid characters", "12345a6789012345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOrderNumberLuna(tt.orderNumber)
			if got != tt.want {
				t.Errorf("ValidateOrderNumberLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
