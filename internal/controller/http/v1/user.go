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

func (server *Server) newUserRoutes(handler *gin.Engine, userUsecase usecase.User, logger logger.Interface) {
	routes := &userRoutes{userUsecase, logger}

	handler.POST("/users", routes.createUser)
	handler.POST("/login", routes.loginUser)
	handler.POST("/logout", routes.logout)
	handler.POST("/renew_token", server.renewAccessToken)

	authRoutes := handler.Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/users/my_profile", routes.getMyProfile)
	authRoutes.PATCH("/users/admin", routes.addAdminRole)
	authRoutes.PATCH("/users/", routes.updateUser)
	authRoutes.PATCH("/users/phone_number/", routes.addPhone)
	authRoutes.DELETE("/users/", routes.deleteUser)
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// @Summary     Create User
// @Description Create new User
// @ID          create-user
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body CreateUserRequest true "register"
// @Success     200 {object} entity.User
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /users/ [post]
func (r *userRoutes) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - createUser")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.CreateUser(&entity.UserRegister{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - createUser")
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

// @Summary     Login
// @Description Log in
// @ID          login
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body LoginUserRequest true "log in"
// @Success     200 {object} entity.UserLoginResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /login/ [post]
func (r *userRoutes) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - loginUser")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.LoginUser(&entity.UserLogin{Email: req.Email, Password: req.Password})
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - loginUser")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.SetCookie("refresh_token", user.RefreshToken, int(time.Until(user.RefreshTokenExpiresAt).Seconds()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, user)
}

// @Summary     Get my profile
// @Description getMyProfile
// @ID          getMyProfile
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.User
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /users/my_profile [get]
func (r *userRoutes) getMyProfile(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	user, st, err := r.userUsecase.GetMyProfile(payload.UserId)
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getMyProfile")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary     Add adminRole
// @Description addAdminRole
// @ID          addAdminRole
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /users/admin [patch]
func (r *userRoutes) addAdminRole(ctx *gin.Context) {
	payload := getJWTPayload(ctx)

	user, st, err := r.userUsecase.AddAdminRole(payload.UserId)
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// @Summary     Update user
// @Description updateUser
// @ID          updateUser
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body UpdateUserRequest true "updateUser"
// @Success     200 {object} entity.User
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /users/ [patch]
func (r *userRoutes) updateUser(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - updateUser")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.UpdateUser(payload.UserId, &entity.UserUpdate{
		Name: req.Name,
	})
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - updateUser")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type AddPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary     AddPhone
// @Description addPhone
// @ID          addPhone
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body AddPhoneRequest true "addPhone"
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /users/phone_number [patch]
func (r *userRoutes) addPhone(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	var req AddPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - addPhone")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp, st, err := r.userUsecase.AddPhone(payload.UserId, &entity.UserAddPhone{
		Phone: req.Phone,
	})
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - addPhone")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary     Delete User
// @Description deleteUser
// @ID          deleteUser
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /users/ [delete]
func (r *userRoutes) deleteUser(ctx *gin.Context) {
	payload := getJWTPayload(ctx)
	res, st, err := r.userUsecase.DeleteUser(payload.UserId)
	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - deleteUser")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary     Logout
// @Description logout
// @ID          logout
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /logout [post]
func (r *userRoutes) logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "logged out")
}
