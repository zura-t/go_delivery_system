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

// @Summary     RenewAccessToken
// @Description renewAccessToken
// @ID          renewAccessToken
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} renewAccessTokenResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /renew_token [post]
func (server *Server) renewAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		server.l.Error(err, "http - v1 - renewAccessToken - context cookie")
		errorResponse(ctx, http.StatusUnauthorized, "can't renew the token")
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		server.l.Error(err, "http - v1 - renewAccessToken - server.tokenMaker.VerifyToken")
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserId, refreshPayload.IsAdmin, refreshPayload.Email, server.config.AccessTokenDuration)
	if err != nil {
		server.l.Error(err, "http - v1 - renewAccessToken - server.tokenMaker.CreateToken")
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
