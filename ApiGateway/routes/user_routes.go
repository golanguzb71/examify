package routes

import (
	"apigateway/handlers"
	"apigateway/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.GET("/profile", middleware.AuthMiddleware([]string{"ROLE_ADMIN", "ROLE_USER"}), handlers.GetUserProfile)
}
