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
	repo *repository.PostgresRepository
}

func NewIeltsService(repo *repository.PostgresRepository) *IeltsService {
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

func (s *IeltsService) UpdateBookById(ctx context.Context, req *pb.UpdateBookRequest) (*pb.AbsResponse, error) {
	err := s.repo.UpdateBook(req.Id, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Book updated successfully",
	}, nil
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

func (s *IeltsService) CreateAttemptInline(ctx context.Context, req *pb.CreateInlineAttemptRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateAttemptInline(req.ExamId, req.UserAnswer, req.SectionType)
	if err != nil {
		log.Printf("Failed to create inline attempt: %v", err)
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Inline attempt created successfully",
	}, nil
}

func (s *IeltsService) CreateAttemptOutlineWriting(ctx context.Context, req *pb.CreateOutlineAttemptRequestWriting) (*pb.AbsResponse, error) {
	err := s.repo.CreateAttemptOutlineWriting(req)
	if err != nil {
		log.Printf("Failed to create outline attempt: %v", err)
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Outline Writing attempt created successfully",
	}, nil
}

func (s *IeltsService) CreateExam(ctx context.Context, req *pb.CreateExamRequest) (*pb.CreateExamResponse, error) {
	examID, err := s.repo.CreateExam(req.UserId, req.BookId)
	if err != nil {
		log.Printf("Failed to create exam: %v", err)
		return nil, err
	}

	return &pb.CreateExamResponse{
		ExamId: *examID,
	}, nil
}

func (s *IeltsService) GetExamByUserId(ctx context.Context, req *pb.GetExamByUserIdRequest) (*pb.GetExamByUserIdResponse, error) {
	return s.repo.GetExamsByUserId(req.UserId, req.PageRequest.Page, req.PageRequest.Size)
}

func (s *IeltsService) GetTopExamResultList(ctx context.Context, req *pb.GetTopExamRequest) (*pb.GetTopExamResult, error) {
	return s.repo.GetTopExamResults(req.Dataframe, req.PageRequest.Page, req.PageRequest.Size)
}

func (s *IeltsService) CreateAttemptOutlineSpeaking(ctx context.Context, req *pb.CreateOutlineAttemptRequestSpeaking) (*pb.AbsResponse, error) {
	err := s.repo.CreateAttemptOutlineSpeaking(req)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "Speaking saved",
	}, nil
}

func (s *IeltsService) GetResultsInlineBySection(ctx context.Context, req *pb.GetResultRequest) (*pb.GetResultResponse, error) {
	return s.repo.GetResultsInlineBySection(req.Section, req.ExamId)
}
