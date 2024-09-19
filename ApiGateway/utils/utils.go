package utils

import "github.com/gin-gonic/gin"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Status  int32       `json:"statusCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"object,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Status  int32  `json:"statusCode"`
	Message string `json:"message"`
}

func RespondSuccess(ctx *gin.Context, statusCode int32, message string, object interface{}) {
	ctx.JSON(int(statusCode), SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    object,
	})
}

func RespondError(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}
