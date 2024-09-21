package handlers

import (
	"apigateway/internal/models"
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateIeltsBook creates a new IELTS book
// @Summary ROLE_ADMIN
// @Description Create a new book for IELTS exam preparation
// @Tags ielts-book
// @Accept  json
// @Produce  json
// @Param name path string true "Name of the book"
// @Success 200 {object} utils.AbsResponse "Book created successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid input"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/book/create/{name} [post]
func CreateIeltsBook(ctx *gin.Context) {
	bookName := ctx.Param("name")
	resp, err := ieltsClient.CreateBook(bookName)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
}

// DeleteIeltsBook deletes an IELTS book by its ID
// @Summary ROLE_ADMIN
// @Description Delete an IELTS book by its ID
// @Tags ielts-book
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the book"
// @Success 200 {object} utils.AbsResponse "Book deleted successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid input"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/book/delete/{id} [delete]
func DeleteIeltsBook(ctx *gin.Context) {
	bookId := ctx.Param("id")
	resp, err := ieltsClient.DeleteBook(bookId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllBook retrieves all IELTS books
// @Summary ROLE_ADMIN
// @Description Retrieve a list of all IELTS books
// @Tags ielts-book
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.AbsResponse "List of books"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/book/books [get]  // Ensure this matches the route definition
func GetAllBook(ctx *gin.Context) {
	books, err := ieltsClient.GetBook()
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// CreateAnswer creates a new answer for an IELTS book
// @Summary ROLE_ADMIN
// @Description Create a new answer for a specified IELTS book
// @Tags ielts-answer
// @Accept  json
// @Produce  json
// @Param bookId path string true "ID of the book"
// @Param answer body models.CreateAnswer true "Answer content"
// @Success 200 {object} utils.AbsResponse "Answer created successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid input"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/answer/create/{bookId} [post]
func CreateAnswer(ctx *gin.Context) {
	bookId := ctx.Param("bookId")
	var req models.CreateAnswer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid input format")
		return
	}
	resp, err := ieltsClient.CreateAnswer(bookId, req.Answers, req.SectionType)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteAnswer deletes an answer by its ID
// @Summary ROLE_ADMIN
// @Description Delete an answer by its ID
// @Tags ielts-answer
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the answer"
// @Success 200 {object} utils.AbsResponse "Answer deleted successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid input"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/answer/delete/{id} [delete]
func DeleteAnswer(ctx *gin.Context) {
	bookId := ctx.Param("id")
	answer, err := ieltsClient.DeleteAnswer(bookId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, answer.Status, answer.Message)
	return
}

// GetAnswerByBookId retrieves an answer by book ID.
//
// @Summary ROLE_ADMIN
// @Description Retrieve the answer associated with a specific book ID via gRPC.
// @Tags ielts-answer
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} pb.GetAnswerResponse "Answer found"
// @Failure 502 {string} string "Error while gRPC connection"
// @Router /api/ielts/answer/{id} [get]
func GetAnswerByBookId(ctx *gin.Context) {
	bookId := ctx.Param("id")
	answer, err := ieltsClient.GetAnswer(bookId)
	if err != nil {
		utils.RespondSuccess(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, answer)
	return
}
