package main

import (
	"authService/config"
	"authService/internal/server"
	"authService/internal/telegram"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	go func() {
		err = server.Run(cfg)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	telegram.RunBot(cfg)
}
