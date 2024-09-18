package client

import (
	"apigateway/proto/pb"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type IeltsClient struct {
	client pb.IeltsServiceClient
}

func NewIeltsClient(addr string) (*IeltsClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewIeltsServiceClient(conn)
	return &IeltsClient{client: client}, nil
}

func (c *IeltsClient) CreateBook(name string) (*pb.CreateBookResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateBookRequest{Name: name}
	resp, err := c.client.CreateBook(ctx, req)
	if err != nil {
		log.Printf("Error when calling CreateBook %v", err)
		return nil, err
	}
	return resp, nil
}
