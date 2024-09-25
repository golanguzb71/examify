package handlers

import (
	"apigateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login godoc
// @Summary ALL
// @Description Validate the authorization code and log in the user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param code path string true "Authorization Code"
// @Success 200 {object} utils.AbsResponse "Login successful"
// @Failure 400 {object} utils.AbsResponse "Invalid code or login failed"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/auth/login/{code} [post]
func Login(ctx *gin.Context) {
	code := ctx.Param("code")
	resp, err := authClient.ValidateCode(code)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, resp.Token)
	return
}
