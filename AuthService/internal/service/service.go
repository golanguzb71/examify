package service

import "authService/proto/pb"

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
