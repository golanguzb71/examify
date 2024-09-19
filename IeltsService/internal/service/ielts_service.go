package service

import (
	"context"
	"ielts-service/proto/pb"
	"log"
)

type BookRepository interface {
	CreateBook(name string) (int32, string, error)
}

type IeltsService struct {
	pb.UnimplementedIeltsServiceServer
	repo BookRepository
}

func NewIeltsService(repo BookRepository) *IeltsService {
	return &IeltsService{repo: repo}
}

func (s *IeltsService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	status, message, err := s.repo.CreateBook(req.Name)
	if err != nil {
		log.Printf("Failed to create book: %v", err)
		return &pb.CreateBookResponse{
			Status:  500,
			Message: "Internal server error",
		}, err
	}

	return &pb.CreateBookResponse{
		Status:  status,
		Message: message,
	}, nil
}
