// Пакет main - главный пакет приложения.
package main

import (
	"log"

	"github.com/RomanAgaltsev/urlcut/internal/app"
)

func main() {
	// Создаем и инициализируем приложение
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize application : %s", err.Error())
	}

	// Запускаем приложение
	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run application : %s", err.Error())
	}
}
