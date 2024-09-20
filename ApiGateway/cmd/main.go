package main

import (
	"apigateway/config"
	_ "apigateway/docs"
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
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
	handlers.InitIeltsClient(grpcClients.IeltsClient)
	handlers.InitUserClient(grpcClients.UserClient)
	handlers.InitAuthClient(grpcClients.AuthClient)
	routes.SetUpRoutes(router)

	port := cfg.Server.Port
	log.Printf("Starting Api Gateway on port %s", port)
	if err = router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type GrpcClients struct {
	IeltsClient *client.IeltsClient
	UserClient  *client.UserClient
	AuthClient  *client.AuthClient
}

func InitializeGrpcClients(cfg *config.Config) *GrpcClients {
	ieltsClient, err := client.NewIeltsClient(cfg.Grpc.IeltsService.Address)
	if err != nil {
		log.Fatalf("Failed to create IELTS client: %v", err)
	}

	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("Failed to create User client: %v", err)
	}

	authClient, err := client.NewAuthClient(cfg.Grpc.AuthService.Address)
	if err != nil {
		return nil
	}
	return &GrpcClients{
		IeltsClient: ieltsClient,
		UserClient:  userClient,
		AuthClient:  authClient,
	}
}
