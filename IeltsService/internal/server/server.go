package server

import (
	"ielts-service/config"
	client "ielts-service/internal/grpc_clients"
	"ielts-service/internal/repository"
	"ielts-service/internal/service"
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

	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("Failed to listen grpc service %v", err)
	}

	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db, userClient)
	ieltsService := service.NewIeltsService(repo)

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
