package main

import (
	"context"
	"log"

	"url/internal/app"
)

func main() {
	ctx := context.Background()

	application := app.NewApp()

	if err := application.Start(ctx); err != nil {
		log.Printf("error while starting application")
	}
}
