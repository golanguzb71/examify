package client

import (
	"context"
	"google.golang.org/grpc"
	"ielts-service/proto/pb"
)

type BonusClient struct {
	client pb.BonusServiceClient
}

func NewBonusClient(addr string) (*BonusClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := pb.NewBonusServiceClient(conn)
	return &BonusClient{client: client}, nil
}

func (c *BonusClient) CreateAttemptBonus(chatId string) bool {
	ctx := context.TODO()
	resp, err := c.client.UseBonusAttempt(ctx, &pb.UseBonusAttemptRequest{ChatId: chatId})
	if err != nil {
		return false
	}
	return resp.Response
}
