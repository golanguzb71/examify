package handlers

import (
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetBonusInformationMe godoc
// @Summary USER
// @Description Retrieve bonus information by the current user's chat ID.
// @Tags bonus
// @Produce json
// @Success 200 {object} pb.GetBonusInformationByChatIdResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 409 {object} map[string]string "Conflict"
// @Router /api/bonus/get-bonus-information-me [get]
func GetBonusInformationMe(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := authClient.GetBonusInformationByChatId(user.ChatId)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
