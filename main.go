package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(Validate("4417123456789113"))
}

func Validate(numberCard string) bool {
	// Алгоритм Луны
	// 1 этап:
	// Создадим слайс для хранения цифр
	var digits []int
	var validCard bool

	// Перебираем каждый символ строки, преобразуем его в число и добавляем в слайс
	for _, char := range numberCard {
		// Преобразуем символ в цифру
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			fmt.Println("Ошибка при преобразовании символа:", err)
			return validCard
		}
		digits = append(digits, digit)
	}
	// 2 этап (умножаем на 2,каждую вторую цифру начиная с первой):
	var digitsMulti []int
	for i, num := range digits {
		if i%2 == 0 {
			if num*2 > 9 { // 3 этап (если в результате удвоения получается две цифры, то сложите эти две цифры):
				sum := (num*2)/10 + (num*2)%10 // num/10 извлекает десятки (например для 35 это будет 3); num%10 извлекает единицы (для 35 это 5)
				digitsMulti = append(digitsMulti, sum)
			} else {
				digitsMulti = append(digitsMulti, num*2)
			}
		} else {
			digitsMulti = append(digitsMulti, num)
		}
	}
	// 4 этап (складываем все числа)
	var sumGrand int // итоговая сумма
	for _, num := range digitsMulti {
		sumGrand += num
	}

	if sumGrand%10 == 0 { // Если итоговая сумма делится на 10, то кредитная карта действительна.
		validCard = true
	} else {
		validCard = false
	}

	return validCard
}
