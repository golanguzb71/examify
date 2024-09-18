package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
