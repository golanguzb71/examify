package handlers

import client "apigateway/internal/grpc_clients"

var (
	ieltsClient *client.IeltsClient
	authClient  *client.AuthClient
	userClient  *client.UserClient
)

func InitIeltsClient(client *client.IeltsClient) {
	ieltsClient = client
}

func InitUserClient(client *client.UserClient) {
	userClient = client
}

func InitAuthClient(client *client.AuthClient) {
	authClient = client
}
