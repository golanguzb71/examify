package service

import (
	"context"
	"user-service/internal/repository"
	"user-service/proto/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserService(repo *repository.PostgresUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.AbsResponse, error) {
	return s.repo.CreateUser(req.Name, req.Surname, req.PhoneNumber, "USER", req.ChatId)
}

func (s *UserService) GetUserByChatIdOrPhone(ctx context.Context, req *pb.GetUserByChatIdOrPhoneRequestOrId) (*pb.User, error) {
	return s.repo.GetUserByChatIdOrPhone(req.ChatId, req.PhoneNumber, req.Id)
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.PageRequest) (*pb.GetAllUserResponse, error) {
	return s.repo.GetAllUsers(req.Page, req.Size)
}

func (s *UserService) UpdateNameSurname(ctx context.Context, req *pb.UpdateUserNameSurnameRequest) (*pb.AbsResponse, error) {
	return s.repo.UpdateNameOrSurname(req.Name, req.Surname)
}
