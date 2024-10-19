package main

import (
	"context"
	"log"

	"github.com/RomanAgaltsev/urlcut/internal/app"
)

func main() {
	ctx := context.Background()
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application : %s", err.Error())
	}
	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run application : %s", err.Error())
	}
}
