package validators

import (
	"regexp"
	"unicode"
)

func ValidateOrderNumber(orderNumber string) bool {
	// Пример простой валидации: номер заказа должен содержать только цифры и быть длиной от 10 до 20 символов
	re := regexp.MustCompile(`^\d{1,20}$`)
	return re.MatchString(orderNumber)
}

func ValidateOrderNumberLuhn(orderNumber string) bool {
	// Удаляем пробелы из номера заказа
	var cleaned string
	for _, r := range orderNumber {
		if unicode.IsDigit(r) {
			cleaned += string(r)
		} else if !unicode.IsSpace(r) {
			return false
		}
	}

	// Проверяем, что номер заказа содержит только цифры
	if len(cleaned) == 0 {
		return false
	}

	// Алгоритм Луна
	var sum int
	length := len(cleaned)
	parity := length % 2
	for i, r := range cleaned {
		digit := int(r - '0')
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
