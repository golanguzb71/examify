package client

import (
	"context"
	"google.golang.org/grpc"
	"ielts-service/proto/pb"
)

type IntegrationClient struct {
	client pb.IntegrationServiceClient
}

func NewIntegrationClient(address string) (*IntegrationClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)))
	if err != nil {
		return nil, err
	}
	client := pb.NewIntegrationServiceClient(conn)
	return &IntegrationClient{client: client}, nil
}

func (c *IntegrationClient) GetResultWritingTask(qua *pb.WritingTaskAbsRequest) (*pb.WritingTaskAbsResponse, error) {
	ctx := context.TODO()
	resp, err := c.client.GetResultWritingTask(ctx, qua)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *IntegrationClient) GetResultSpeakingPart(req *pb.CreateOutlineAttemptRequestSpeaking) (*pb.SpeakingPartAbsResponse, error) {
	resp, err := c.client.GetResultSpeakingPart(context.TODO(), &pb.SpeakingPartAbsRequest{
		Question:     req.GetQuestion(),
		VoiceMessage: req.GetVoiceAnswer(),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
