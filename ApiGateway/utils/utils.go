package utils

import "github.com/gin-gonic/gin"

// AbsResponse represents an error API response
type AbsResponse struct {
	Status  int32  `json:"statusCode"`
	Message string `json:"message"`
}

func RespondSuccess(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), AbsResponse{
		Status:  statusCode,
		Message: message,
	})
}

func RespondError(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), AbsResponse{
		Status:  statusCode,
		Message: message,
	})
}
