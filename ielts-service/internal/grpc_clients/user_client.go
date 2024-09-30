package client

import (
	"context"
	"google.golang.org/grpc"
	"ielts-service/proto/pb"
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

func (c *UserClient) GetUserByPhoneNumberOrChatIdOrId(phoneNumber, chatId, id *string) *pb.User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb.GetUserByChatIdOrPhoneRequestOrId{
		ChatId:      chatId,
		PhoneNumber: phoneNumber,
		Id:          id,
	}
	User, err := c.client.GetUserByChatIdOrPhone(ctx, &req)
	if err != nil {
		return nil
	}
	return User
}
