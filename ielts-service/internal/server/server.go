package server

import (
	"ielts-service/config"
	client "ielts-service/internal/grpc_clients"
	"ielts-service/internal/repository"
	"ielts-service/internal/service"
	"ielts-service/internal/utils"
	"ielts-service/proto/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// grpc Clients making
	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	integrationClient, err := client.NewIntegrationClient(cfg.Grpc.IntegrationService.Address)

	if err != nil {
		log.Fatalf("Failed to listen grpc service %v", err)
	}

	// db connecting
	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// database migration
	utils.MigrateUp(db)

	// repository_making
	repo := repository.NewPostgresRepository(db, userClient, integrationClient)
	// service making
	ieltsService := service.NewIeltsService(repo)

	// grpc server making
	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIeltsServiceServer(grpcServer, ieltsService)

	log.Printf("Server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
