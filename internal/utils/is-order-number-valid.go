package utils

import "strconv"

func IsOrderNumberValid(orderNumber string) bool {
	sum := 0
	parity := len(orderNumber) % 2
	for i, c := range orderNumber {
		digit, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}

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
