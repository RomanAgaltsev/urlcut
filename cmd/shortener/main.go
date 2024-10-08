package main

import (
	"log"

	"github.com/RomanAgaltsev/urlcut/internal/app"
)

func main() {
	// Создаем новое приложение
	a, err := app.New()
	// Проверяем наличие ошибок
	if err != nil {
		// Есть ошибка, выводим
		log.Fatalf("running shortener application failed: %s", err.Error())
	}
	// Запускаем приложение
	a.Run()
}
