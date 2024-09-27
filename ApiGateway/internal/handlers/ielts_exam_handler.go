package handlers

import (
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetExamResult(c *gin.Context) {

}

func GetExamUserAnswers(ctx *gin.Context) {

}

// GetTopExamResult Retrieve top exam results based on the specified dataframe (MONTHLY, DAILY, or WEEKLY)
// @Summary ALL
// @Description Retrieve top exam results based on the specified dataframe (MONTHLY, DAILY, or WEEKLY)
// @Tags ielts-exam
// @Accept json
// @Produce json
// @Param dataframe path string true "The timeframe for which to retrieve top exam results (MONTHLY, DAILY, WEEKLY)"
// @Param page query int true "The page number for pagination"
// @Param size query int true "The number of results per page"
// @Success 200 {object} pb.GetTopExamResult "Successful response with exam results"
// @Failure 400 {object} utils.AbsResponse "Bad request with error message"
// @Router /api/ielts/exam/result/top-exam-result/{dataframe} [get]
func GetTopExamResult(ctx *gin.Context) {
	value := ctx.Param("dataframe")
	p := ctx.Query("page")
	s := ctx.Query("size")
	page, err := strconv.Atoi(p)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "error while parsing page")
		return
	}
	size, err := strconv.Atoi(s)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "error while parsing size")
		return
	}
	if value != "MONTHLY" && value != "DAILY" && value != "WEEKLY" {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid dataframe")
		return
	}
	result, err := ieltsClient.GetTopExamResult(value, page, size)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

// CreateExam godoc
// @Summary USER
// @Description This endpoint creates a new exam for the specified user and book.
// @Tags ielts-exam
// @Accept  json
// @Produce  json
// @Param  bookId  query  int  true  "Book ID"
// @Success 200 {object} utils.AbsResponse "Exam created successfully, returning the exam ID"
// @Failure 400 {object} utils.AbsResponse "Invalid input parameters"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/ielts/exam/create [post]
// @Security Bearer
func CreateExam(ctx *gin.Context) {
	bookIdStr := ctx.Query("bookId")
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bookId, err := strconv.ParseInt(bookIdStr, 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	examId, err := ieltsClient.CreateExam(int32(user.Id), int32(bookId))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, examId.ExamId)
	return
}
