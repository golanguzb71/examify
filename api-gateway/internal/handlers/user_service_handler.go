package handlers

import (
	"apigateway/proto/pb"
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetUserProfile godoc
// @Summary USER
// @Description Retrieves the profile information of the currently authenticated user
// @Tags user
// @Produce  json
// @Success 200 {object} pb.User "User Profile Response"
// @Failure 409 {object} utils.AbsResponse "Conflict Error"
// @Router /api/user/profile [get]
// @Security Bearer
func GetUserProfile(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	response, _ := authClient.CalculateBonusToday(user.ChatId)
	remainCount := ieltsClient.CalculateTodayExamCount(user.Id)
	user.TodayExamCount = response + remainCount
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// ChangeNameSurname godoc
// @Summary USER
// @Description Allows a user to update their name and surname
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body pb.UpdateUserNameSurnameRequest true "Name and Surname Update Request"
// @Success 200 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse "Conflict Error"
// @Router /api/user/update-information [put]
// @Security Bearer
func ChangeNameSurname(ctx *gin.Context) {
	var req pb.UpdateUserNameSurnameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	req.UserId = strconv.FormatInt(user.Id, 10)
	resp, err := userClient.UpdateNameSurname(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetMyResults godoc
// @Summary USER
// @Description Retrieves a paginated list of exam results for the logged-in user
// @Tags ielts-exam
// @Accept  json
// @Produce  json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Success 200 {object} pb.GetExamByUserIdResponse "User's exam results"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 401 {object} utils.AbsResponse "Unauthorized"
// @Failure 502 {object} utils.AbsResponse "Bad Gateway"
// @Router /api/user/get-my-results [get]
// @Security Bearer
func GetMyResults(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	page, err := strconv.ParseInt(ctx.Query("page"), 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	size, err := strconv.ParseInt(ctx.Query("size"), 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := ieltsClient.GetExamByUserId(int32(user.Id), int32(page), int32(size))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
