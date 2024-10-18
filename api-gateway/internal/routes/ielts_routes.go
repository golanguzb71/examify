package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func IELTSRoutes(r *gin.RouterGroup, authClient *client.AuthClient) {
	ielts := r.Group("/ielts")

	book := ielts.Group("/book")
	{
		book.POST("/create/:name", middleware.AuthMiddleware([]string{"ADMIN"}, authClient), handlers.CreateIeltsBook)
		book.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN"}, authClient), handlers.DeleteIeltsBook)
		book.PUT("/update/:id", middleware.AuthMiddleware([]string{"ADMIN"}, authClient), handlers.UpdateIeltsBook)
		book.GET("/books", handlers.GetAllBook)
	}
	answer := ielts.Group("/answer", middleware.AuthMiddleware([]string{"ADMIN"}, authClient))
	{
		answer.POST("/create/:bookId", handlers.CreateAnswer)
		answer.DELETE("/delete/:answerId", handlers.DeleteAnswer)
		answer.PUT("/update/:id", handlers.UpdateAnswer)
		answer.GET("/:id", handlers.GetAnswerByBookId)
	}

	exam := ielts.Group("/exam")
	{
		exam.POST("/create", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.CreateExam)
		result := exam.Group("/result")
		{
			result.GET("/:page/:size", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetExamResult)
			result.GET("/user-answers/:examId", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetExamUserAnswers)
			result.GET("/top-exam-result/:dataframe", handlers.GetTopExamResult)
			result.GET("/get-results-inline/:sectionType/:examId", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetResultsInlineBySection)
			result.GET("/get-results-outline-writing/:examId", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetResultsOutlineWriting)
			result.GET("/get-results-outline-speaking/:examId/:partNumber", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetResultsOutlineSpeaking)
			result.GET("/get-results-speaking-voice/:voiceName", handlers.GetVoiceRecord)
		}
		attempt := exam.Group("/attempt")
		{
			attempt.POST("/create/inline", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.CreateInlineAttempt)
			attempt.POST("/create/outline-writing", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.CreateOutlineAttemptWriting)
			attempt.POST("/create/outline-speaking/:examId", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.CreateOutlineAttemptSpeaking)
		}

	}
}
