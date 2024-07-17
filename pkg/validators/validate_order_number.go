package validators

import (
	"regexp"
	"strconv"
	"strings"
)

func ValidateOrderNumber(orderNumber string) bool {
	// Пример простой валидации: номер заказа должен содержать только цифры и быть длиной от 10 до 20 символов
	re := regexp.MustCompile(`^\d{10,20}$`)
	return re.MatchString(orderNumber)
}

func ValidateOrderNumberLuna(orderNumber string) bool {
	// Удаляем пробелы из номера заказа
	orderNumber = strings.ReplaceAll(orderNumber, " ", "")

	// Проверяем, что номер заказа содержит только цифры
	_, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return false
	}

	// Алгоритм Луна
	var sum int
	length := len(orderNumber)
	parity := length % 2
	for i, digit := range orderNumber {
		digit := int(digit - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
