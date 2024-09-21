package service

import (
	"context"
	"ielts-service/internal/repository"
	"ielts-service/proto/pb"
	"log"
	"net/http"
)

type IeltsService struct {
	pb.UnimplementedIeltsServiceServer
	repo *repository.PostgresBookRepository
}

func NewIeltsService(repo *repository.PostgresBookRepository) *IeltsService {
	return &IeltsService{repo: repo}
}

func (s *IeltsService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateBook(req.Name)
	if err != nil {
		log.Printf("Failed to create book: %v", err)
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "Book created successfully",
	}, nil
}

func (s *IeltsService) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteBook(req.BookId)
	if err != nil {
		log.Printf("Failed to delete book: %v", err)
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "Book deleted successfully",
	}, nil
}

func (s *IeltsService) GetAllBook(ctx context.Context, req *pb.Empty) (*pb.GetAllBookResponse, error) {
	books, err := s.repo.GetAllBooks()
	if err != nil {
		log.Printf("Failed to get all books: %v", err)
		return nil, err
	}

	var protoBooks []*pb.Book
	for _, book := range books {
		protoBooks = append(protoBooks, &pb.Book{
			Id:    book.ID,
			Title: book.Title,
		})
	}

	return &pb.GetAllBookResponse{Books: protoBooks}, nil
}

func (s *IeltsService) CreateAnswer(ctx context.Context, req *pb.CreateAnswerRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateAnswer(req.BookId, req.SectionType, req.Answers)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "Answer created successfully",
	}, nil
}

func (s *IeltsService) DeleteAnswer(ctx context.Context, req *pb.DeleteBookRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteAnswer(req.BookId)
	if err != nil {
		log.Printf("Failed to delete answer: %v", err)
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Answer deleted successfully",
	}, nil
}

func (s *IeltsService) GetAnswer(ctx context.Context, req *pb.GetAnswerRequest) (*pb.GetAnswerResponse, error) {
	answers, err := s.repo.GetAnswerByBookId(req.BookId)
	if err != nil {
		log.Printf("Failed to get answers for book ID: %v", err)
		return nil, err
	}

	var protoAnswers []*pb.Answer
	for _, answer := range answers {
		protoAnswers = append(protoAnswers, &pb.Answer{
			Id:            answer.ID,
			BookId:        answer.BookId,
			SectionType:   answer.SectionType,
			SectionAnswer: answer.Answer,
		})
	}

	return &pb.GetAnswerResponse{Answers: protoAnswers}, nil
}
