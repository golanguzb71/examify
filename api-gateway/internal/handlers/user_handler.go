package handlers

import (
	"apigateway/proto/pb"
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserProfile godoc
// @Summary Get user profile
// @Description Retrieves the profile information of the currently authenticated user
// @Tags User
// @Produce  json
// @Success 200 {object} pb.User "User Profile Response"
// @Failure 409 {object} utils.AbsResponse "Conflict Error"
// @Router /api/user/profile [get]
// @Security Bearer
func GetUserProfile(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// ChangeNameSurname godoc
// @Summary Update user's name and surname
// @Description Allows a user to update their name and surname
// @Tags user-default
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
	resp, err := userClient.UpdateNameSurname(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
}
