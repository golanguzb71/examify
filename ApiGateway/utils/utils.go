package utils

import (
	"github.com/gin-gonic/gin"
)

func RespondSuccess(ctx *gin.Context, statusCode int32, message string, object interface{}) {
	ctx.JSON(int(statusCode), gin.H{
		"statusCode": statusCode,
		"message":    message,
		"object":     object,
	})
}

func RespondError(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), gin.H{
		"statusCode": statusCode,
		"message":    message,
		"object":     nil,
	})
}
