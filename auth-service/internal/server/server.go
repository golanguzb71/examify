package server

import (
	"authService/config"
	client "authService/internal/grpc_clients"
	database "authService/internal/repository"
	"authService/internal/service"
	"authService/internal/telegram"
	"authService/proto/pb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(cfg *config.Config) error {
	database.ConnectRedis(cfg)
	userClient, err := client.NewUserClient(cfg.GRPC.UserService.Address)
	if err != nil {
		return err
	}
	telegram.SetUserClient(userClient)
	authService := service.NewAuthService(userClient)
	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Server listening on port %s", cfg.Server.Port)
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
