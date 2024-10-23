package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserServiceRoutes(r *gin.RouterGroup, authClient *client.AuthClient) {
	user := r.Group("/user")
	user.GET("/profile", middleware.AuthMiddleware([]string{"ADMIN", "USER"}, authClient), handlers.GetUserProfile)
	user.PUT("/update-information", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.ChangeNameSurname)
	user.GET("/get-my-results", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetMyResults)
}
