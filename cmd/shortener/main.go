// Пакет main - главный пакет приложения.
package main

import (
	"log"

	"github.com/RomanAgaltsev/urlcut/internal/app"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Выводим информацию о сборке
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)

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
