package logic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InfArgStdin(useIssuers bool, useBrands bool, once bool, stdinInput bool, inf bool) bool {
	flags2 := make(map[string]func(string) func(string))
	args := os.Args
	if stdinInput && inf {
		// Если флаг --stdin активирован, читаем номера карт с stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Читаем строку с номерами карт
			inputLine := scanner.Text()
			if inputLine == "" {
				fmt.Fprintf(os.Stderr, "input line not found\n")
			}

			// Разбиваем строку на отдельные номера карт
			cardNumbers := strings.Fields(inputLine)

			// Обрабатываем каждый номер карты
			for _, cardNumber := range cardNumbers {
				if Validate(cardNumber) == false {
					fmt.Println(cardNumber)
					fmt.Println("Correct: no")
					fmt.Println("Card Brand: -")
					fmt.Println("Card Issuer: -")
					continue // Если карта некорректна, то прекращаем дальнейшую обработку для этой карты
				}
				once = true
				for i := 1; i < len(args); i++ {
					if args[i] == "information" {
						// Перебираем флаги --brands и --issuers
						if useIssuers && strings.HasPrefix(args[i+1], "--issuers=") {
							i++
						}
						if strings.HasPrefix(args[i+1], "--brands=") {
							brandsFile := strings.TrimPrefix(args[i+1], "--brands=")
							flags2["information"](brandsFile)(cardNumber)
						} else if args[i+1] == "--brands" && i+1 < len(args) {
							flags2["information"](args[i+1])(cardNumber)
						}
						if useBrands && strings.HasPrefix(args[i+1], "--brands=") {
							i++
						}
						if useIssuers && strings.HasPrefix(args[i-1], "--issuers=") {
							i--
							i--
						}
						if strings.HasPrefix(args[i+1], "--issuers=") {
							issuersFile := strings.TrimPrefix(args[i+1], "--issuers=")
							flags2["information"](issuersFile)(cardNumber)
						} else if args[i+1] == "--issuers" && i+1 < len(args) {
							flags2["information"](args[i+1])(cardNumber)
						}
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
	} else {
		// Обработка аргументов командной строки
		for i := 1; i < len(args); i++ {
			// Обработка флага "information"
			if args[i] == "information" && i+1 < len(args) {
				// Перебираем номера карт после флага "information"
				for j := i + 1; j < len(args); j++ {
					cardNumber := args[j] // Номер карты берется из аргумента

					if cardNumber == "--brands=brands.txt" || cardNumber == "--issuers=issuers.txt" {
						continue
					}
					if Validate(cardNumber) == false {
						fmt.Println(cardNumber)
						fmt.Println("Correct: no")
						fmt.Println("Card Brand: -")
						fmt.Println("Card Issuer: -")
						continue // Если карта некорректна, то прекращаем дальнейшую обработку для этой карты
					}
					once = true
					if useBrands || useIssuers {
						// Если флаги активированы, передаем в функцию
						for k := i + 1; k < len(args); k++ {
							if strings.HasPrefix(args[k], "--brands=") {
								brandsFile := strings.TrimPrefix(args[k], "--brands=")
								flags2["information"](brandsFile)(cardNumber)
							}
							if strings.HasPrefix(args[k], "--issuers=") {
								issuersFile := strings.TrimPrefix(args[k], "--issuers=")
								flags2["information"](issuersFile)(cardNumber)
							}
						}
					}
				}
			}
		}
	}
	return once
}
