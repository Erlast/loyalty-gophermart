package validators

import (
	"regexp"
)

func ValidateOrderNumber(orderNumber string) bool {
	// Пример простой валидации: номер заказа должен содержать только цифры и быть длиной от 10 до 20 символов
	re := regexp.MustCompile(`^\d{10,20}$`)
	return re.MatchString(orderNumber)
}
