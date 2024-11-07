package logic

import (
	"math/rand"
	"time"
)

func randomInt(min, max int) int {
	//
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min // Генерируем случайное число в диапазоне [min, max)
}
