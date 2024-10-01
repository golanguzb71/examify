package client

import (
	"apigateway/proto/pb"
	"context"
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

func (uc *UserClient) UpdateNameSurname(req *pb.UpdateUserNameSurnameRequest) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return uc.client.UpdateNameSurname(ctx, req)
}
