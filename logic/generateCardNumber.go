package logic

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Функция для генерации номера карты
func GenerateCardNumber(prefix string, length int) string {
	var cardNumber string
	for Validate(cardNumber) == false { // До тех пор пока номер карты не валиднен
		// Устанавливаем начальное зерно для генератора случайных чисел
		// Это необходимо для того, чтобы каждый запуск программы генерировал разные случайные числа
		rand.Seed(time.Now().UnixNano())

		// Инициализация строки cardNumber, начиная с префикса
		cardNumber = prefix

		// Генерация случайных цифр до достижения нужной длины (должна быть длина минус 1, т.к. последняя цифра - контрольная)
		// В цикле добавляются случайные цифры (от 0 до 9) в cardNumber
		for len(cardNumber) < length-1 {
			cardNumber += fmt.Sprintf("%d", rand.Intn(10)) // Добавляем случайное число от 0 до 9
		}

		// Алгоритм Луна для вычисления контрольной цифры:
		// Этот алгоритм проверяет и рассчитывает контрольную цифру для проверки корректности номера карты.
		// 1 этап:
		// Создадим слайс для хранения цифр
		var digits []int

		// Перебираем каждый символ строки, преобразуем его в число и добавляем в слайс
		for i := len(cardNumber) - 1; i >= 0; i-- {
			// Преобразуем символ в цифру
			digit, err := strconv.Atoi(string(cardNumber[i]))
			if err != nil {
				fmt.Println("Symbol does not convert to number")
				os.Exit(1)
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

		// Рассчитываем контрольную цифру
		// Сначала находим остаток от деления суммы на 10 (это будет последняя цифра суммы)
		checksum := (10 - sumGrand%10) % 10

		// Добавляем контрольную цифру в конец номера карты
		cardNumber += fmt.Sprintf("%d", checksum)

	}
	// Возвращаем готовый номер карты
	return cardNumber
}
