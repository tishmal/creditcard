package logic

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// функция для флага issue
func Issue(brand string, issuer string, useBrands bool, useIssuers bool, brandsFile string, issuersFile string) {
	// Если бренд и эмитент равны пустоте, то выход со статусом 1
	if brand == "" && issuer == "" {
		fmt.Println("brand and issuer not found")
		os.Exit(1)
	}
	// Если эмитент пуст и флаг эмитента активирован, то выход со статусом 1
	if issuer == "" && useIssuers {
		fmt.Println("flag --issuers use, but --issuer not found")
		os.Exit(1)
	}
	// Если бренд равен пустоте и флаг бренда актвирован, то выход со статусом 1
	if brand == "" && useBrands {
		fmt.Println("flag --brands use, but --brand not found")
		os.Exit(1)
	}
	// Если используется флаг --brands или --issuers, обрабатываем их
	if useBrands || useIssuers { // код который ниже как оказалось подходит логически и практически для обработки данных с обоих текстовых файлов, поэтому чтобы не дублировать код, создадим метод и будем прогонять два вида данных
		// brand - полиморфизм, может быть brand или issuers в зависимости от того какой файл пришёл в параметры флага information
		// Читаем файл с брендами
		brandsFileContent, err := os.Open(brandsFile)
		issuersFileContent, err2 := os.Open(issuersFile)
		if err != nil && err2 != nil {
			fmt.Println("information: file .txt not found")
			os.Exit(1)
		}
		defer brandsFileContent.Close()
		defer issuersFileContent.Close()

		// Проверяем, пустой ли файл
		var stat fs.FileInfo
		var stat2 fs.FileInfo
		if useBrands {
			stat, err = brandsFileContent.Stat()
			if err != nil {
				fmt.Println("information: could not get file stats")
				os.Exit(1)
			}
			if stat.Size() == 0 {
				fmt.Println("information: file is empty")
				os.Exit(1)
			}
		}
		if useIssuers {
			stat2, err2 = issuersFileContent.Stat()
			if err2 != nil {
				fmt.Println("information: could not get file stats")
				os.Exit(1)
			}
			if stat2.Size() == 0 {
				fmt.Println("information: file is empty")
				os.Exit(1)
			}
		}
		// Чтение файла brandsFile построчно
		scanner := bufio.NewScanner(brandsFileContent)
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
		// Чтение файла issuersFile построчно
		scanner2 := bufio.NewScanner(issuersFileContent)
		if err := scanner2.Err(); err != nil {
			os.Exit(1)
		}

		var iin string
		var bin string
		if useIssuers {
			iin = ReturnIINorBIN(scanner2, issuer, useBrands, useIssuers, iin)
		}
		if useBrands {
			bin = ReturnIINorBIN(scanner, brand, useBrands, useIssuers, iin)
		}
		if strings.HasPrefix(iin, bin) {
			// Генерация номера карты
			fmt.Println(GenerateCardNumber(iin, 16))
		} else if useBrands && !useIssuers {
			fmt.Println(GenerateCardNumber(bin, 16))
		}
	} else {
		os.Exit(1)
	}
}
