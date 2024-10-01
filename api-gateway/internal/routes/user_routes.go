package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, authClient *client.AuthClient) {
	user := r.Group("/user")
	user.GET("/profile", middleware.AuthMiddleware([]string{"ADMIN", "USER"}, authClient), handlers.GetUserProfile)
}