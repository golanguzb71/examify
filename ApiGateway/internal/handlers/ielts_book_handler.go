package handlers

import (
	"apigateway/internal/grpc_clients"
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ieltsClient *client.IeltsClient

func InitIeltsClient(client *client.IeltsClient) {
	ieltsClient = client
}

// CreateIeltsBookRequest represents the request body for creating a book
// @Summary Create a new IELTS book
// @Description Create a new book for IELTS
// @Tags IELTS
// @Accept  json
// @Produce  json
// @Param name path string true "Name of the book"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/ielts/book/create/{name} [post]
func CreateIeltsBook(ctx *gin.Context) {
	bookName := ctx.Param("name")
	resp, err := ieltsClient.CreateBook(bookName)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message, nil)
}

func DeleteIeltsBook(ctx *gin.Context) {

}

func GetAllBook(ctx *gin.Context) {

}

func CreateAnswer(ctx *gin.Context) {

}

func DeleteAnswer(ctx *gin.Context) {

}
