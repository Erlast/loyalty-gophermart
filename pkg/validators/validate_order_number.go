package validators

import (
	"unicode"
)

func ValidateOrderNumber(orderNumber string) bool {
	// Удаляем пробелы и проверяем на нецифровые символы
	var cleaned string
	for _, r := range orderNumber {
		if unicode.IsDigit(r) {
			cleaned += string(r)
		} else if !unicode.IsSpace(r) {
			// Возвращаем false при обнаружении любого нецифрового и не пробельного символа
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
