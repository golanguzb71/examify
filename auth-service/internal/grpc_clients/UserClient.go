package client

import (
	"authService/proto/pb"
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client}, nil
}

func (c *UserClient) GetUserByChatIdOrPhone(chatID, phoneNumber, id *string) (*pb.User, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("UserClient or gRPC client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetUserByChatIdOrPhoneRequestOrId{
		ChatId:      chatID,
		PhoneNumber: phoneNumber,
		Id:          id,
	}
	return c.client.GetUserByChatIdOrPhone(ctx, req)
}

func (c *UserClient) CreateUser(name, surname, chatID, phoneNumber string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateUserRequest{
		Name:        name,
		Surname:     surname,
		ChatId:      chatID,
		PhoneNumber: phoneNumber,
	}
	return c.client.CreateUser(ctx, req)
}
