package routes

import (
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.GET("/profile", middleware.AuthMiddleware([]string{"ROLE_ADMIN", "ROLE_USER"}), handlers.GetUserProfile)
}
