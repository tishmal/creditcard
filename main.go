package main

import (
	"bufio"
	"creditcard/logic"
	"fmt"
	"os"
	"strings"
)

// // Функция флага information
// func useBrandsAndOrIssuers(_cardNumber string, _numBrand string) bool {
// }

func main() {
	args := os.Args

	if len(args) < 3 {
		os.Exit(1)
	}
	flags := make(map[string]func(string))
	flags2 := make(map[string]func(string) func(string))

	var useBrands bool   // используется ли флаг --brands
	var useIssuers bool  // используется ли флаг --issuers
	var once bool = true // повторить один раз
	var inf bool
	isPick := args[2] == "--pick"

	flags2["information"] = func(brandsFile string) func(cardNumber string) {
		return func(cardNumber string) {
			// Если используется флаг --brands или --issuers, обрабатываем их
			if useBrands || useIssuers { // код который ниже как оказалось подходит логически и практически для обработки данных с обоих текстовых файлов, поэтому чтобы не дублировать код, создадим метод и будем прогонять два вида данных
				// brand - полиморфизм, может быть brand или issuers в зависимости от того какой файл пришёл в параметры флага information
				var hasBrand bool    // имеет ли карта бренд схожий с брендом в текстовом файле
				var nameBrand string // имя бренда
				var numBrand string  // номер бренда
				// Читаем файл с брендами
				brandsFileContent, err := os.Open(brandsFile)
				if err != nil {
					fmt.Println("information: file .txt not found")
					os.Exit(1)
				}
				defer brandsFileContent.Close()

				// Чтение файла построчно
				scanner := bufio.NewScanner(brandsFileContent)
				if err := scanner.Err(); err != nil {
					os.Exit(1)
				}
				for scanner.Scan() {
					line := scanner.Text() // тут каждую строку текстового файла записываем в line. line это например VISA:4
					// Разделяем строку на бренд и номер
					for i, ch := range line { // перебираем строку line и если встретим в ней ':' то дробим строку
						if ch == ':' {
							nameBrand = line[:i]
							numBrand = line[i+1:]
						}
					}
					if strings.HasPrefix(cardNumber, numBrand) { // если карта имеет префикс бренда, то...
						hasBrand = true // карта имеет бренд
						break           // программа выходит из цикла
					}
				}
				var _nameIssuers string // имя эмитента
				var _nameBr string      // имя бренда

				if useIssuers && hasBrand {
					_nameIssuers = nameBrand
				}
				if useBrands && hasBrand {
					_nameBr = nameBrand
				}
				// Печатаем информацию для карты
				if once == true {
					fmt.Println(cardNumber)
					fmt.Println("Correct: yes")
					once = false
				}

				if useIssuers && brandsFile == "issuers.txt" {
					fmt.Println("Card Issuer:", _nameIssuers)
				}
				if useBrands && brandsFile == "brands.txt" {
					fmt.Println("Card Brand:", _nameBr)
				}
			}
		}
	}

	flags["generate"] = func(value string) { // value - это значение(номер кредитной карты или нет) которое мы получаем из консоли
		logic.GenerateCards(value, isPick)
	}

	var hasError bool

	flags["validate"] = func(value string) {
		if logic.Validate(value) == false {
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
	// Обработка флагов --brands и --issuers
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
	if stdinInput && inf {
		// Если флаг --stdin активирован, читаем номера карт с stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Читаем строку с номерами карт
			inputLine := scanner.Text()
			if inputLine == "" {
				fmt.Fprintf(os.Stderr, "INCORRECT\n")
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
			fmt.Fprintf(os.Stderr, "INCORRECT\n")
			hasError = true
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
					if logic.Validate(cardNumber) == false {
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

	if stdinInput {
		// Если флаг --stdin активирован, читаем номера карт с stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Читаем одну строку, содержащую все номера карт
			inputLine := scanner.Text()
			if inputLine == "" {
				fmt.Fprintf(os.Stderr, "INCORRECT\n")
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
				if isPick == false {
					flags["generate"](args[i+1])
				}
				if isPick == true {
					if len(os.Args) < 4 {
						os.Exit(1)
					}
					flags["generate"](args[i+2])
				}
				return
			}
		}
	}
	if hasError {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
