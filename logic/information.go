package logic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Information(emissioFile string, cardNumber string, useBrands bool, useIssuers bool, once bool) bool {
	// Если используется флаг --brands или --issuers, обрабатываем их
	if useBrands && useIssuers { // код который ниже как оказалось подходит логически и практически для обработки данных с обоих текстовых файлов, поэтому чтобы не дублировать код, создадим метод и будем прогонять два вида данных
		// brand - полиморфизм, может быть brand или issuers в зависимости от того какой файл пришёл в параметры флага information
		var hasBrand bool    // имеет ли карта бренд схожий с брендом в текстовом файле
		var nameBrand string // имя бренда
		var numBrand string  // номер бренда
		// Читаем файл с брендами
		brandsFileContent, err := os.Open(emissioFile)
		if err != nil {
			fmt.Println("information: file .txt not found")
			os.Exit(1)
		}
		defer brandsFileContent.Close()

		// Проверяем, пустой ли файл
		stat, err := brandsFileContent.Stat()
		if err != nil {
			fmt.Println("information: could not get file stats")
			os.Exit(1)
		}
		if stat.Size() == 0 {
			fmt.Println("information: file is empty")
			os.Exit(1)
		}

		// Чтение файла построчно
		scanner := bufio.NewScanner(brandsFileContent)
		if err := scanner.Err(); err != nil {
			fmt.Println("Reading file with error")
			os.Exit(1)
		}
		for scanner.Scan() {
			line := scanner.Text() // тут каждую строку текстового файла записываем в line. line это например VISA:4
			if line == "" {
				fmt.Println("Reading file with error")
				os.Exit(1)
			}
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
		if emissioFile == "brands.txt" {
			if _nameBr == "" {
				_nameBr = "-"
			}
			fmt.Println("Card Brand:", _nameBr)
		}
		if emissioFile == "issuers.txt" {
			if _nameIssuers == "" {
				_nameIssuers = "-"
			}
			fmt.Println("Card Issuer:", _nameIssuers)
		}
	} else {
		fmt.Println("flags --brands or --issuers not found")
		os.Exit(1)
	}
	return once
}
