package routes

import (
	client "apigateway/internal/grpc_clients"
	"apigateway/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(r *gin.Engine, authClient *client.AuthClient) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/auth/login/:code", handlers.Login)
		UserRoutes(api, authClient)
		IELTSRoutes(api, authClient)
	}
}
