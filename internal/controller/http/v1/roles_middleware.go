package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) rolesMiddleware() gin.HandlerFunc {
	abort := func(ctx *gin.Context, err error) {
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		ctx.Abort()
	}

	return func(ctx *gin.Context) {
		jwtPayload := getJWTPayload(ctx)
		user, _, err := server.userUsecase.GetMyProfile(jwtPayload.UserId)

		if err != nil {
			abort(ctx, errors.New("Can't get payload"))
			return
		}

		if !user.IsAdmin {
			abort(ctx, errors.New("Incorrect user role"))
			return
		}

		ctx.Next()
	}
}
