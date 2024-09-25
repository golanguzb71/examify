package service

import (
	client "authService/internal/grpc_clients"
	"authService/internal/security"
	"authService/internal/telegram"
	"authService/proto/pb"
	"context"
	"errors"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userClient *client.UserClient
}

func NewAuthService(userServiceClient *client.UserClient) *AuthService {
	return &AuthService{
		userClient: userServiceClient,
	}
}

func (s *AuthService) ValidateCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ValidateCodeResponse, error) {
	chatId := telegram.GetStoredCode(req.Code)
	if chatId == nil {
		return nil, errors.New("invalid code")
	}
	user, err := s.userClient.GetUserByChatIdOrPhone(chatId, nil, nil)
	if err != nil {
		return nil, err
	}
	token, err := security.GenerateToken(user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateCodeResponse{User: user, Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.User, error) {
	claims, err := security.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	user, err := s.userClient.GetUserByChatIdOrPhone(nil, &claims.Username, nil)
	if err != nil {
		return nil, err
	}
	var checker = false
	for _, role := range req.RequiredRoles {
		if user.Role == role {
			checker = true
			break
		}
	}
	if !checker {
		return nil, errors.New("this user not valid for this endpoint")
	}
	return user, nil
}
