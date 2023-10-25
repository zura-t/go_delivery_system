package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/internal/entity"
	"github.com/zura-t/go_delivery_system/internal/usecase"
	"github.com/zura-t/go_delivery_system/pkg/logger"
)

type userRoutes struct {
	userUsecase usecase.User
	logger      logger.Interface
}

func (server *Server) newUserRoutes(handler *gin.RouterGroup, userUsecase usecase.User, logger logger.Interface) {
	routes := &userRoutes{userUsecase, logger}

	handler.POST("/users", routes.createUser)
	handler.POST("/login", routes.loginUser)
	handler.POST("/logout", routes.logout)
	handler.POST("/renew_token", server.renewAccessToken)

	authRoutes := handler.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/users/my_profile", routes.getMyProfile)
	authRoutes.PATCH("/users/", routes.updateUser)
	authRoutes.PATCH("/users/phone_number/", routes.addPhone)
	authRoutes.DELETE("/users/", routes.deleteUser)
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *userRoutes) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.CreateUser(entity.UserRegister{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	User                  entity.User `json:"user"`
}

func (r *userRoutes) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.LoginUser(entity.UserLogin{Email: req.Email, Password: req.Password})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.SetCookie("refresh_token", user.RefreshToken, int(time.Until(user.RefreshTokenExpiresAt).Seconds()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, user)
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *userRoutes) getMyProfile(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	user, st, err := r.userUsecase.GetMyProfile(payload.UserId)
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (r *userRoutes) updateUser(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.UpdateUser(payload.UserId, entity.UserUpdate{
		Name: req.Name,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type AddPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (r *userRoutes) addPhone(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	var req AddPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp, st, err := r.userUsecase.AddPhone(payload.UserId, entity.UserAddPhone{
		Phone: req.Phone,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (r *userRoutes) deleteUser(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	res, st, err := r.userUsecase.DeleteUser(payload.UserId)
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (r *userRoutes) logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "logged out")
}
