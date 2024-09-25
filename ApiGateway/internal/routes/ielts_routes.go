package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func IELTSRoutes(r *gin.RouterGroup, authClient *client.AuthClient) {
	ielts := r.Group("/ielts")

	book := ielts.Group("/book", middleware.AuthMiddleware([]string{"ADMIN"}, authClient))
	{
		book.POST("/create/:name", handlers.CreateIeltsBook)
		book.DELETE("/delete/:id", handlers.DeleteIeltsBook)
		book.GET("/books", handlers.GetAllBook)
	}

	answer := ielts.Group("/answer", middleware.AuthMiddleware([]string{"ADMIN"}, authClient))
	{
		answer.POST("/create/:bookId", handlers.CreateAnswer)
		answer.DELETE("/delete/:answerId", handlers.DeleteAnswer)
		answer.GET("/:id", handlers.GetAnswerByBookId)
	}

	exam := ielts.Group("/exam")
	{
		result := exam.Group("/result")
		{
			result.GET("/:page/:size", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetExamResult)
			result.GET("/user-answers/:examId", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetExamUserAnswers)
			result.GET("/top-exam-result/:dataframe", handlers.GetTopExamResult)
		}
	}
}
