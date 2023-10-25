package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/token"
)

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "can't renew the token")
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserId, refreshPayload.Email, server.config.AccessTokenDuration)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "can't create new token")
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func getJWTPayload(ctx *gin.Context) token.Payload {
	var payload token.Payload
	payloadData, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		errorResponse(ctx, http.StatusInternalServerError, "couldn't get payload from authtoken")
		return token.Payload{}
	}
	data, ok := payloadData.(token.Payload)
	if ok {
		payload = data
	} else {
		errorResponse(ctx, http.StatusInternalServerError, "couldn't get payload from authtoken")
		return token.Payload{}
	}
	return payload
}