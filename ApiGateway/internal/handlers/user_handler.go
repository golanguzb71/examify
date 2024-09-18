package handlers

import (
	"apigateway/internal/grpc_clients"
	"github.com/gin-gonic/gin"
)

var userClient *client.UserClient

func InitUserClient(client *client.UserClient) {
	userClient = client
}

func GetUserProfile(ctx *gin.Context) {

}
