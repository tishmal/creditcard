package main

import (
	"bufio"
	"creditcard/logic"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 3 { // Если длина аргументов меньше 3, то завершаем программу
		fmt.Println("Incorrect input. Arguments not recognized.")
		os.Exit(1)
	}
	flags := make(map[string]func(string))
	flags3 := make(map[string]func(string) func(string) func(string))
	flags4 := make(map[string]func(string) func(string) func(string) func(string))

	var useBrands bool            // Используется ли флаг --brands
	var useIssuers bool           // Используется ли флаг --issuers
	var once bool = true          // Повторить один раз
	var inf bool                  // Используется ли флаг --information. необходим для активации --stdin и -information вместе
	isPick := args[2] == "--pick" // Используется ли флаг --pick

	// Поиск флагов --brands и --issuers
	for i := 1; i < len(args); i++ {
		if args[i] == "--brands" || strings.HasPrefix(args[i], "--brands=") {
			useBrands = true
		}
		if args[i] == "--issuers" || strings.HasPrefix(args[i], "--issuers=") {
			useIssuers = true
		}
	}
	for i := 1; i < len(args); i++ {
		if args[i] == "information" {
			inf = true
		}
	}
	// Поиск флага --stdin
	var stdinInput bool
	for i := 1; i < len(args); i++ {
		if args[i] == "--stdin" {
			stdinInput = true
		}
	}
	// Флаги:
	flags4["issue"] = func(brand string) func(issuer string) func(brandsFile string) func(issuersFile string) {
		// обрабатываем параметры после обработки командной строки и обработки флагов текстовых документов
		return func(issuer string) func(brandsFile string) func(issuersFile string) {
			return func(brandsFile string) func(issuersFile string) {
				return func(issuersFile string) {
					logic.Issue(brand, issuer, useBrands, useIssuers, brandsFile, issuersFile)
				}
			}
		}
	}

	flags3["information"] = func(brandsFile string) func(issuersFile string) func(cardNumber string) {
		return func(issuersFile string) func(cardNumber string) {
			return func(cardNumber string) {
				if logic.Validate(cardNumber) == false {
					fmt.Println("number card invalid")
					os.Exit(1)
				}
				infa := logic.Information(brandsFile, cardNumber, useBrands, useIssuers, once)
				if infa == false {
					once = false
				} else {
					once = true
				}
				infa2 := logic.Information(issuersFile, cardNumber, useBrands, useIssuers, once)
				if infa2 == false {
					once = true
				} else {
					once = true
				}
			}
		}
	}

	flags["generate"] = func(value string) { // value - это значение(номер кредитной карты или нет) которое мы получаем из консоли
		for i := 0; i < len(args); i++ {
			if i > 3 && isPick {
				fmt.Println("Unnecessary arguments")
				os.Exit(1)
			}
			if i > 2 && !isPick {
				fmt.Println("Unnecessary arguments")
				os.Exit(1)
			}
		}
		logic.GenerateCards(value, isPick)
	}

	flags["validate"] = func(value string) {
		if logic.Validate(value) == false {
			fmt.Fprintf(os.Stderr, "INCORRECT\n")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stdout, "OK\n")
		}
	}

	// Обработка флагов и аргументов
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
				if logic.Validate(cardNumber) == false {
					fmt.Println(cardNumber)
					fmt.Println("Correct: no")
					fmt.Println("Card Brand: -")
					fmt.Println("Card Issuer: -")
					continue // Если карта некорректна, то прекращаем дальнейшую обработку для этой карты
				}

				once = true
				var brandsFile string
				var issuersFile string
				for i := 1; i < len(args); i++ {
					if args[i] == "information" {
						// Перебираем флаги --brands и --issuers
						if useIssuers && strings.HasPrefix(args[i+1], "--issuers=") {
							i++
						}
						if strings.HasPrefix(args[i+1], "--brands=") {
							brandsFile = strings.TrimPrefix(args[i+1], "--brands=")
						}
						if useBrands && strings.HasPrefix(args[i+1], "--brands=") {
							i++
						}
						if useIssuers && strings.HasPrefix(args[i-1], "--issuers=") {
							i--
							i--
						}
						if strings.HasPrefix(args[i+1], "--issuers=") {
							issuersFile = strings.TrimPrefix(args[i+1], "--issuers=")
						}
					}
				}
				flags3["information"](brandsFile)(issuersFile)(cardNumber)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scaner")
			os.Exit(1)
		}
	} else {
		// Обработка аргументов командной строки
		for i := 1; i < len(args); i++ {
			// Обработка флага "information"
			if args[i] == "information" && i+1 < len(args) {
				// Перебираем номера карт после флага "information"
				for j := i + 3; j < len(args); j++ {
					cardNumber := args[j] // Номер карты берется из аргумента

					if cardNumber == "--brands=brands.txt" || cardNumber == "--issuers=issuers.txt" {
						continue
					}
					once = true
					if useBrands || useIssuers {
						var brandsFile string
						var issuersFile string
						// Если флаги активированы, передаем в функцию
						for k := i + 1; k < len(args); k++ {
							if strings.HasPrefix(args[k], "--brands=") {
								brandsFile = strings.TrimPrefix(args[k], "--brands=")
							}
							if strings.HasPrefix(args[k], "--issuers=") {
								issuersFile = strings.TrimPrefix(args[k], "--issuers=")
							}
						}
						brandsFileContent, err := os.Open(brandsFile)
						issuersFileContent, err2 := os.Open(issuersFile)
						if err != nil || err2 != nil {
							fmt.Println("information: file .txt not found")
							os.Exit(1)
						}
						defer brandsFileContent.Close()
						defer issuersFileContent.Close()

						if logic.Validate(cardNumber) == false {
							fmt.Println(cardNumber)
							fmt.Println("Correct: no")
							fmt.Println("Card Brand: -")
							fmt.Println("Card Issuer: -")
							continue // Если карта некорректна, то прекращаем дальнейшую обработку для этой карты
						}
						flags3["information"](brandsFile)(issuersFile)(cardNumber)
					}
				}
			}
		}
	}

	if stdinInput {
		for i := 0; i < len(args); i++ {
			if args[len(args)-1] != "--stdin" {
				fmt.Println("Flag --stdin not found")
				os.Exit(1)
			}
		}
		// Если флаг --stdin активирован, читаем номера карт с stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Читаем одну строку, содержащую все номера карт
			inputLine := scanner.Text()
			if inputLine == "" {
				fmt.Fprintf(os.Stderr, "INCORRECT\n")
				os.Exit(1)
			}

			// Разбиваем строку на отдельные номера карт
			cardNumbers := strings.Fields(inputLine)

			// Обрабатываем каждый номер карты
			for _, cardNumber := range cardNumbers {
				flags["validate"](cardNumber) // Проверяем каждый номер карты
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "INCORRECT\n")
			os.Exit(1)
		}
	} else {
		// Обработка аргументов командной строки и вызов флага
		for i := 1; i < len(args); i++ {
			if args[i] == "validate" && i+1 < len(args) {
				// Проверяем все аргументы после ключа "validate"
				for j := i + 1; j < len(args); j++ {
					flags["validate"](args[j]) // Проверяем каждый номер карты
				}
				return // Завершаем выполнение программы после обработки всех карт
			}
			if args[i] == "generate" && i+1 < len(args) {
				if isPick == false {
					flags["generate"](args[i+1])
				}
				if isPick == true {
					if len(os.Args) < 4 {
						fmt.Println("incorrect input")
						os.Exit(1)
					}
					flags["generate"](args[i+2])
				}
				return
			}
			// Обработка аргументов командной строки с использованием флага issue
			if args[i] == "issue" && i+1 < len(args) {
				var brand string                     // бренд карты
				var issuer string                    // эмитент карты
				var brandsFile string                // текстовый файл с брендами
				var issuersFile string               // текстовый файл с эмитентами
				for k := i + 1; k < len(args); k++ { // перебираем аргументы ком. строки
					if strings.HasPrefix(args[k], "--brand=") {
						brand = strings.TrimPrefix(args[k], "--brand=")
					}
					if strings.HasPrefix(args[k], "--issuer=") {
						issuer = strings.TrimPrefix(args[k], "--issuer=")
					}

					if strings.HasPrefix(args[k], "--brands=") {
						brandsFile = strings.TrimPrefix(args[k], "--brands=")
					}
					if strings.HasPrefix(args[k], "--issuers=") {
						issuersFile = strings.TrimPrefix(args[k], "--issuers=")
					}
				}
				flags4["issue"](brand)(issuer)(brandsFile)(issuersFile) // вызов флага
			}
		}
	}
	if !isPick && !useBrands && !useIssuers {
		fmt.Println("Invalid input")
		os.Exit(1)
	}
}
