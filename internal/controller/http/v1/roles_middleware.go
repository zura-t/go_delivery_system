package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/token"
)

func rolesMiddleware() gin.HandlerFunc {
	abort := func(ctx *gin.Context, err error) {
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		ctx.Abort()
	}

	return func(ctx *gin.Context) {
		var tokenPayload token.Payload
		payload, exists := ctx.Get(authorizationPayloadKey)
		tokenPayload = payload.(token.Payload)

		if exists {
			abort(ctx, errors.New("Can't get payload"))
			return
		}

		if !tokenPayload.IsAdmin {
			abort(ctx, errors.New("Incorrect user role"))
			return
		}

		ctx.Next()
	}
}
