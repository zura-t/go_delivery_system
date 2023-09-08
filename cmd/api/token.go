package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserId, refreshPayload.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
