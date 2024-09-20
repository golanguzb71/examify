package routes

import (
	"apigateway/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app := r.Group("/api")
	app.POST("/auth/login/:code", handlers.Login)
	UserRoutes(app)
	IELTSRoutes(app)
}
