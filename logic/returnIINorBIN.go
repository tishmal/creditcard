package logic

import (
	"bufio"
	"os"
	"strings"
)

// Функция для флага issue, чтобы не дублировать код, она возвращает идентификационный номер бренда или эмитента (IIN)
// emissio это brand или issuer карты
func ReturnIINorBIN(scanner *bufio.Scanner, emissio string, useBrands bool, useIssuers bool, iin string) string {
	var name string  // имя бренда или эмитента
	var num string   // номер бренда или эмитента
	var count int    // счётчик, сколько имеется одинаковых брендов или эмитентов в текстовом файле
	var name1 string // это имя бренда или эмитента который сравнивается с входящим именем и если совпадает, то счётчик прибавляется

	// Сохраняем все строки в срез
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Проходим по сохранённым строкам для подсчёта совпадений с эмитентом
	for _, line1 := range lines {
		for i, ch := range line1 {
			if ch == ':' {
				name1 = line1[:i]
			}
		}
		if name1 == emissio {
			count++
		}
	}

	// Второй проход по сохранённым строкам для основной логики
	for _, line := range lines {
		if line == "" {
			os.Exit(1)
		}

		// Разделяем строку на бренд и номер
		for i, ch := range line {
			if ch == ':' {
				name = line[:i]
				num = line[i+1:]
			}
		}

		// Логика с учетом флагов useBrands и useIssuers
		if useBrands && useIssuers {
			if strings.HasPrefix(iin, num) == false && count > 1 {
				// Мы уже прошли все строки, так что тут можно просто перебирать lines снова
				for _, line := range lines {
					for i, ch := range line {
						if ch == ':' {
							name = line[:i]
							num = line[i+1:]
						}
					}
					if strings.HasPrefix(iin, num) {
						return num
					}
				}
				return num
			}
		}

		// Основная логика поиска
		if name == emissio {
			return num
		}
	}

	// Если ничего не найдено
	return ""
}
