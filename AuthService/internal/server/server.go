package server

import (
	"authService/config"
	"authService/internal/service"
	"authService/proto/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ieltsService := service.NewAuthService()

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, ieltsService)

	log.Printf("Server listening on port %s", cfg.Server.Port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
