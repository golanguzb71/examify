package server

import (
	"database/sql"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"ielts-service/config"
	client "ielts-service/internal/grpc_clients"
	"ielts-service/internal/repository"
	"ielts-service/internal/service"
	"ielts-service/proto/pb"
	"log"
	"net"
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
	bonusClient, err := client.NewBonusClient(cfg.Grpc.BonusService.Address)
	if err != nil {
		log.Fatalf("Failed to listen grpc service %v", err)
		return
	}
	// db connecting
	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// database migration
	//utils.MigrateUp(db)

	// repository_making
	repo := repository.NewPostgresRepository(db, userClient, integrationClient, bonusClient)
	// service making
	ieltsService := service.NewIeltsService(repo)

	c := cron.New()
	_, err = c.AddFunc("@every 1m", func() {
		if err := updatePendingExamsStatus(db); err != nil {
			log.Printf("Error updating pending exams status: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to create cron job: %v", err)
	}

	c.Start()
	defer c.Stop()

	// grpc server making
	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10 MB
		grpc.MaxSendMsgSize(10*1024*1024), // 10 MB
	)
	pb.RegisterIeltsServiceServer(grpcServer, ieltsService)

	log.Printf("Server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func updatePendingExamsStatus(db *sql.DB) error {
	_, err := db.Exec("SELECT update_pending_exams_status();")
	if err != nil {
		log.Printf("Error executing SQL function: %v", err)
		return err
	}
	return nil
}
