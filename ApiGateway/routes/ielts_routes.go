package routes

import (
	"apigateway/handlers"
	"apigateway/middleware"
	"github.com/gin-gonic/gin"
)

func IELTSRoutes(r *gin.RouterGroup) {
	ielts := r.Group("/ielts")

	book := ielts.Group("/book", middleware.AuthMiddleware([]string{"ROLE_ADMIN"}))
	book.POST("/create/:name", handlers.CreateIeltsBook)
	book.DELETE("/delete/:id", handlers.DeleteIeltsBook)
	book.GET("/getAllBook", handlers.GetAllBook)

	answer := book.Group("/answer", middleware.AuthMiddleware([]string{"ROLE_ADMIN"}))
	answer.POST("/create/:bookId", handlers.CreateAnswer)
	answer.DELETE("/delete/answerId", handlers.DeleteAnswer)

	exam := ielts.Group("/exam")
	result := exam.Group("/result")
	result.GET("/:page/:size", middleware.AuthMiddleware([]string{"ROLE_USER"}), handlers.GetExamResult)
	result.GET("/user-answers/:examId", middleware.AuthMiddleware([]string{"ROLE_USER"}), handlers.GetExamUserAnswers)
	result.GET("/top-exam-result/:dataframe}", handlers.GetTopExamResult)

}
