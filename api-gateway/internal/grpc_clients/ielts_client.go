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

func (c *IeltsClient) CreateExam(userId, bookId int32) (*pb.CreateExamResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateExamRequest{
		UserId: userId,
		BookId: bookId,
	}
	response, err := c.client.CreateExam(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *IeltsClient) GetExamByUserId(userId, page, size int32) (*pb.GetExamByUserIdResponse, error) {
	return nil, nil
}

func (c *IeltsClient) CreateAttemptInline(CAI *pb.CreateInlineAttemptRequest) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.client.CreateAttemptInline(ctx, CAI)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *IeltsClient) CreateAttemptOutlineWriting(coaw *pb.CreateOutlineAttemptRequestWriting) (*pb.AbsResponse, error) {
	ctx := context.TODO()
	resp, err := c.client.CreateAttemptOutlineWriting(ctx, coaw)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *IeltsClient) GetTopExamResult(dataframe string, page, size int) (*pb.GetTopExamResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.GetTopExamRequest{
		Dataframe: dataframe,
		PageRequest: &pb.PageRequest{
			Page: int32(page),
			Size: int32(size),
		},
	}
	list, err := c.client.GetTopExamResultList(ctx, req)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *IeltsClient) UpdateBookById(id, name string) (*pb.AbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.UpdateBookRequest{
		Id:   id,
		Name: name,
	}
	resp, err := c.client.UpdateBookById(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
