package logic

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GenerateCards(value string, isPick bool) {
	var countStar int
	var lastValueThisStar bool

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
		var validNumbers []string
		for i, char := range value {
			if i >= len(value)-countStar {
				maxCombinations := 1
				for i := 0; i < countStar; i++ {
					maxCombinations *= 10
				}

				for i := 0; i < maxCombinations; i++ {
					// Форматируем число с ведущими нулями
					replacement := fmt.Sprintf("%0*d", countStar, i)
					// Создаем текущий номер, заменяя звездочки на цифры
					currentNumber := value[:len(value)-countStar] + replacement

					if Validate(currentNumber) {
						validNumbers = append(validNumbers, currentNumber)
					}
				}
				if len(validNumbers) == 0 {
					os.Exit(1)
				}

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
		if isPick {
			rand.Seed(time.Now().UnixNano())
			fmt.Println(validNumbers[rand.Intn(len(validNumbers))])
		} else {
			for _, num := range validNumbers {
				fmt.Println(num)
			}
		}
	}
}
