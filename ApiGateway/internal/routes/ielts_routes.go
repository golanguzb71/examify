package routes

import (
	handlers2 "apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func IELTSRoutes(r *gin.RouterGroup) {
	ielts := r.Group("/ielts")

	book := ielts.Group("/book", middleware.AuthMiddleware([]string{"ROLE_ADMIN"}))
	book.POST("/create/:name", handlers2.CreateIeltsBook)
	book.DELETE("/delete/:id", handlers2.DeleteIeltsBook)
	book.GET("/getAllBook", handlers2.GetAllBook)

	answer := book.Group("/answer", middleware.AuthMiddleware([]string{"ROLE_ADMIN"}))
	answer.POST("/create/:bookId", handlers2.CreateAnswer)
	answer.DELETE("/delete/answerId", handlers2.DeleteAnswer)

	exam := ielts.Group("/exam")
	result := exam.Group("/result")
	result.GET("/:page/:size", middleware.AuthMiddleware([]string{"ROLE_USER"}), handlers2.GetExamResult)
	result.GET("/user-answers/:examId", middleware.AuthMiddleware([]string{"ROLE_USER"}), handlers2.GetExamUserAnswers)
	result.GET("/top-exam-result/:dataframe}", handlers2.GetTopExamResult)

}
