package middleware

import (
	client "apigateway/internal/grpc_clients"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(requiredRoles []string, authClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, bearerPrefix)

		user, err := authClient.ValidateToken(token, requiredRoles)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or insufficient permissions"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
