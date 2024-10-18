package server

import (
	"google.golang.org/grpc"
	"integration-service/config"
	"integration-service/internal/service"
	"integration-service/proto/pb"
	"log"
	"net"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("%v", err.Error())
	}
	integrationService := service.NewIntegrationService()
	listen, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen port %v ,  %v", cfg.Server.Port, err)
	}
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10 MB
		grpc.MaxSendMsgSize(10*1024*1024), // 10 MB
	)
	pb.RegisterIntegrationServiceServer(grpcServer, integrationService)
	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
