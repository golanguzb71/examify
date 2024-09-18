package main

import (
	"apigateway/config"
	_ "apigateway/docs"
	"apigateway/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load confg: %v", err)
	}

	routes.SetUpRoutes(router)
	port := cfg.Server.Port
	log.Printf("Starting Api Gateway on port %s", port)
	if err = router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server:  %v", err)
	}
}
