package main

import (
	"apigateway/config"
	_ "apigateway/docs"
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// @title Examify Swagger

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	router.Use(cors.New(corsConfig))

	grpcClients := InitializeGrpcClients(cfg)
	handlers.InitIeltsClient(grpcClients.IeltsClient)
	handlers.InitUserClient(grpcClients.UserClient)
	handlers.InitAuthClient(grpcClients.AuthClient)
	routes.SetUpRoutes(router, grpcClients.AuthClient)

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
