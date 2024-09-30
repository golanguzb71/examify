package service

import (
	"context"
	"integration-service/proto/pb"
)

type IntegrationService struct {
	pb.UnimplementedIntegrationServiceServer
}

func NewIntegrationService() *IntegrationService {
	return &IntegrationService{}
}

func (i *IntegrationService) GetResultWritingTask(ctx context.Context, wta *pb.WritingTaskAbsRequest) (*pb.WritingTaskAbsResponse, error) {
	return processEssay(wta.Question + " " + wta.Answer)
}

func (i *IntegrationService) GetResultSpeakingPart(ctx context.Context, spr *pb.SpeakingPartAbsRequest) (*pb.SpeakingPartAbsResponse, error) {
	return nil, nil
}
