package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args

	flags := make(map[string]func(string))

	var countStar int
	var lastValueThisStar bool
	flags["generate"] = func(value string) { // value - это значение(номер кредитной карты или нет) которое мы получаем из консоли
		for _, ch := range value {
			if ch == '*' {
				countStar++
			}
			if len(value)-countStar >= 0 && len(value)-countStar < len(value) { // чтобы индекс не вышел из диапазона надо проверить, чтобы он был меньше длины входящей строки value
				if value[len(value)-countStar] == '*' {
					lastValueThisStar = true
				} else {
					lastValueThisStar = false
				}
			}
			if countStar > 4 {
				os.Exit(1)
			}
		}
		if lastValueThisStar == false {
			os.Exit(1)
		}
		if lastValueThisStar == true {
			var digits []int
			for i, char := range value {
				if i >= len(value)-countStar {
					// Генерируем случайное число от 0 до 9
					randomNumber := randomInt(0, 10)

					digits = append(digits, randomNumber)
				} else {
					// Преобразуем символ в цифру
					digit, err := strconv.Atoi(string(char))
					if err != nil {
						return
					}
					digits = append(digits, digit)
				}
			}
			var numCard string
			for _, d := range digits {
				numCard += strconv.Itoa(d)
			}
			fmt.Println(numCard)
		}
	}

	var hasError bool

	flags["validate"] = func(value string) {
		if Validate(value) == false {
			fmt.Fprintf(os.Stderr, "INCORRECT\n")
			hasError = true
		} else {
			fmt.Fprintf(os.Stdout, "OK\n")
		}
	}
	// Обработка флагов и аргументов
	var stdinInput bool
	for i := 1; i < len(args); i++ {
		if args[i] == "--stdin" {
			stdinInput = true
		}
	}

	if stdinInput {
		// Если флаг --stdin активирован, читаем номера карт с stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Читаем строки с картами и передаем их в функцию валидации
			cardNumber := scanner.Text()
			flags["validate"](cardNumber)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "INCORRECT\n")
			hasError = true
		}
	} else {
		// Обработка аргументов командной строки
		for i := 1; i < len(args); i++ {
			if args[i] == "validate" && i+1 < len(args) {
				// Проверяем все аргументы после ключа "validate"
				for j := i + 1; j < len(args); j++ {
					flags["validate"](args[j]) // Проверяем каждый номер карты
				}
				return // Завершаем выполнение программы после обработки всех карт
			}
			if args[i] == "generate" && i+1 < len(args) {
				flags["generate"](args[2])
				return
			}
		}
	}
	// после завершения всех проверок - конец  завершения
	if hasError {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func Validate(numberCard string) bool {
	// Алгоритм Луны
	// 1 этап:
	// Создадим слайс для хранения цифр
	var digits []int
	var validCard bool

	if len(numberCard) < 13 || len(numberCard) > 19 {
		validCard = false
	}

	// Перебираем каждый символ строки, преобразуем его в число и добавляем в слайс
	for _, char := range numberCard {
		// Преобразуем символ в цифру
		digit, err := strconv.Atoi(string(char))
		if err != nil {
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
	if numberCard == "" {
		validCard = false
	}

	return validCard
}

func randomInt(min, max int) int {
	//
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min // Генерируем случайное число в диапазоне [min, max)
}
