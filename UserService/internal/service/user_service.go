package service

import (
	"context"
	"log"
	"user-service/proto/pb"
)

type UserRepository interface {
	CreateUser(name, surname, phoneNumber, role string) (int32, string, error)
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	status, message, err := s.repo.CreateUser(req.Name, req.Surname, req.PhoneNumber, req.Role)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return &pb.CreateUserResponse{
			Status:  500,
			Message: "Internal server error",
		}, err
	}

	return &pb.CreateUserResponse{
		Status:  status,
		Message: message,
	}, nil
}
