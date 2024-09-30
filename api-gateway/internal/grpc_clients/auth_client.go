package client

import (
	"apigateway/proto/pb"
	"context"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{client: client}, nil
}

func (c *AuthClient) ValidateCode(code string) (*pb.ValidateCodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ValidateCodeRequest{Code: code}
	resp, err := c.client.ValidateCode(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *AuthClient) ValidateToken(token string, requiredRoles []string) (*pb.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	req := &pb.ValidateTokenRequest{
		Token:         token,
		RequiredRoles: requiredRoles,
	}
	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
