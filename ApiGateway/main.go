package main

import (
	"apigateway/config"
	client2 "apigateway/internal/grpc_clients"
	handlers2 "apigateway/internal/handlers"
	"apigateway/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	grpcClients := InitializeGrpcClients(cfg)
	handlers2.InitIeltsClient(grpcClients.IeltsClient)
	handlers2.InitUserClient(grpcClients.UserClient)

	routes.SetUpRoutes(router)

	port := cfg.Server.Port
	log.Printf("Starting Api Gateway on port %s", port)
	if err = router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type GrpcClients struct {
	IeltsClient *client2.IeltsClient
	UserClient  *client2.UserClient
}

func InitializeGrpcClients(cfg *config.Config) *GrpcClients {
	ieltsClient, err := client2.NewIeltsClient(cfg.Grpc.IeltsService.Address)
	if err != nil {
		log.Fatalf("Failed to create IELTS client: %v", err)
	}

	userClient, err := client2.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("Failed to create User client: %v", err)
	}

	return &GrpcClients{
		IeltsClient: ieltsClient,
		UserClient:  userClient,
	}
}
