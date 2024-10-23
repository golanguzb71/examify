package client

import (
	"apigateway/proto/pb"
	"context"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	authClient  pb.AuthServiceClient
	bonusClient pb.BonusServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	authClient := pb.NewAuthServiceClient(conn)
	bonusClient := pb.NewBonusServiceClient(conn)
	return &AuthClient{authClient: authClient, bonusClient: bonusClient}, nil
}

func (c *AuthClient) ValidateCode(code string) (*pb.ValidateCodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := &pb.ValidateCodeRequest{Code: code}
	resp, err := c.authClient.ValidateCode(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *AuthClient) ValidateToken(token string, requiredRoles []string) (*pb.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	req := &pb.ValidateTokenRequest{
		Token:         token,
		RequiredRoles: requiredRoles,
	}
	resp, err := c.authClient.ValidateToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *AuthClient) CalculateBonusToday(chatId string) (int32, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	req := &pb.BonusServiceAbsRequest{ChatId: chatId}
	response, err := c.bonusClient.CalculateBonusByChatId(ctx, req)
	if err != nil {
		return 0, err
	}
	return response.Count, nil
}

func (c *AuthClient) GetBonusInformationByChatId(chatId string) (*pb.GetBonusInformationByChatIdResponse, error) {
	return c.bonusClient.GetBonusInformationByChatId(context.TODO(), &pb.BonusServiceAbsRequest{ChatId: chatId})
}
