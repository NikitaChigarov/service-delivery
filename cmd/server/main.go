package main

import (
	"log"
	"service-delivery/internal/app"
	"service-delivery/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
