package handlers

import (
	"apigateway/proto/pb"
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

// CreateInlineAttempt godoc
// @Summary USER
// @Description Creates a new inline attempt for IELTS
// @Tags attempts
// @Accept json
// @Produce json
// @Param request body pb.CreateInlineAttemptRequest true "Create inline attempt request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/attempt/create/inline [post]
func CreateInlineAttempt(ctx *gin.Context) {
	var createInlineAttemptRequest *pb.CreateInlineAttemptRequest
	err := ctx.ShouldBindJSON(&createInlineAttemptRequest)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "json body cannot parsed")
		return
	}
	resp, err := ieltsClient.CreateAttemptInline(createInlineAttemptRequest)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// CreateOutlineAttemptWriting godoc
// @Summary USER
// @Description Creates a new inline attempt for IELTS
// @Tags attempts
// @Accept json
// @Produce json
// @Param request body pb.CreateOutlineAttemptRequestWriting true "Create outline attempt request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/attempt/create/outline-writing [post]
func CreateOutlineAttemptWriting(ctx *gin.Context) {
	var createOutlineAttemptRequestWriting *pb.CreateOutlineAttemptRequestWriting
	err := ctx.ShouldBindJSON(&createOutlineAttemptRequestWriting)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	resp, err := ieltsClient.CreateAttemptOutlineWriting(createOutlineAttemptRequestWriting)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
}

// CreateOutlineAttemptSpeaking godoc
// @Summary USER
// @Description Creates a new inline attempt for IELTS
// @Tags attempts
// @Accept multipart/form-data
// @Produce json
// @Param examId path string true "Exam ID"
// @Param question query string true "Question"
// @Param partNumber query string true "Part Number"
// @Param voiceAnswer formData file true "Voice Answer (only wav)"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/attempt/create/outline-speaking/{examId} [post]
func CreateOutlineAttemptSpeaking(ctx *gin.Context) {
	strExamId := ctx.Param("examId")
	strQuestion := ctx.Query("question")
	file, err := ctx.FormFile("voiceAnswer")
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Voice answer file is required")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".wav" {
		utils.RespondError(ctx, http.StatusBadRequest, "Only only wav files are allowed")
		return
	}
	fileContent, err := file.Open()
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Could not open the uploaded file")
		return
	}
	defer fileContent.Close()
	buffer := make([]byte, 512)
	_, err = fileContent.Read(buffer)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Could not read the file content")
		return
	}

	fileContent.Seek(0, 0)

	strPartNumber := ctx.Query("partNumber")
	partNumber, err := strconv.ParseInt(strPartNumber, 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid part number")
		return
	}

	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Could not read the file content")
		return
	}

	var req pb.CreateOutlineAttemptRequestSpeaking
	req.ExamId = strExamId
	req.Question = strQuestion
	req.VoiceAnswer = fileBytes
	req.PartNumber = int32(partNumber)

	resp, err := ieltsClient.CreateOutlineSpeakingAttempt(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetResultsInlineBySection godoc
// @Summary Get results by section
// @Description Retrieves the inline results of a specific section in an IELTS exam
// @Tags results
// @Accept json
// @Produce json
// @Param sectionType path string true "Section Type"
// @Param examId path string true "Exam ID"
// @Success 200 {object} pb.GetResultResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/result/get-results-inline/{sectionType}/{examId} [get]
func GetResultsInlineBySection(ctx *gin.Context) {
	sectionType := ctx.Param("sectionType")
	examId := ctx.Param("examId")
	response, err := ieltsClient.GetResultsInlineBySection(sectionType, examId)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// GetResultsOutlineWriting godoc
// @Summary Get outline of Writing results
// @Description Retrieves the outline results of the Writing section in an IELTS exam
// @Tags results
// @Accept json
// @Produce json
// @Param examId path string true "Exam ID"
// @Success 200 {object} pb.GetResultOutlineWritingResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/result/get-results-outline-writing/{examId} [get]
func GetResultsOutlineWriting(ctx *gin.Context) {
	examId := ctx.Param("examId")
	response, err := ieltsClient.GetResultsOutlineWriting(examId)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// GetResultsOutlineSpeaking godoc
// @Summary Get outline of Speaking results
// @Description Retrieves the outline results of the Speaking section in an IELTS exam
// @Tags results
// @Accept json
// @Produce json
// @Param examId path string true "Exam ID"
// @Success 200 {object} pb.GetResultOutlineSpeakingResponse
// @Failure 409 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/ielts/exam/result/get-results-outline-speaking/{examId} [get]
func GetResultsOutlineSpeaking(ctx *gin.Context) {
	examId := ctx.Param("examId")
	response, err := ieltsClient.GetResultsOutlineSpeaking(examId)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}
