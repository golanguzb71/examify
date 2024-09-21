package client

import (
	"apigateway/internal/models"
	"apigateway/proto/pb"
	"context"
	"google.golang.org/grpc"
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

func (c *IeltsClient) CreateBook(name string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateBookRequest{Name: name}
	resp, err := c.client.CreateBook(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *IeltsClient) DeleteBook(bookId string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.DeleteBookRequest{BookId: bookId}
	resp, err := c.client.DeleteBook(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *IeltsClient) GetBook() ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.client.GetAllBook(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	var books []*models.Book
	for _, b := range response.Books {
		books = append(books, &models.Book{
			ID:    b.Id,
			Title: b.Title,
		})
	}
	return books, nil
}

func (c *IeltsClient) CreateAnswer(bookId string, answers []string, section string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateAnswerRequest{BookId: bookId, Answers: answers, SectionType: section}
	resp, err := c.client.CreateAnswer(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *IeltsClient) DeleteAnswer(bookId string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DeleteBookRequest{BookId: bookId}
	resp, err := c.client.DeleteAnswer(ctx, req)
	if err != nil {
		return &pb.AbsResponse{
			Status:  resp.Status,
			Message: resp.Message,
		}, err
	}

	return &pb.AbsResponse{
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}

func (c *IeltsClient) GetAnswer(bookId string) (*pb.GetAnswerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.GetAnswerRequest{BookId: bookId}
	answers, err := c.client.GetAnswer(ctx, req)
	if err != nil {
		return nil, err
	}
	return answers, nil
}
