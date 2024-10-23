package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"apigateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func BonusServiceRoutes(api *gin.RouterGroup, authClient *client.AuthClient) {
	bonus := api.Group("/bonus")
	{
		bonus.GET("/get-bonus-information-me", middleware.AuthMiddleware([]string{"USER"}, authClient), handlers.GetBonusInformationMe)
	}
}
