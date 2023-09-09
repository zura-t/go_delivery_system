package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/pb"
	"github.com/zura-t/go_delivery_system/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type userResponse struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	user, err := c.CreateUser(context, &pb.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := c.LoginUser(context, &pb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.SetCookie("refresh_token", res.RefreshToken, int(res.RefreshTokenExpiresAt.AsTime().Sub(time.Now())), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) getMyProfile(ctx *gin.Context) {
	var payload token.Payload
	log.Println(ctx)
	payloadData, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		error := fmt.Errorf("couldn't get payload from authtoken")
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	data, ok := payloadData.(token.Payload)
	if ok {
		payload = data
	} else {
		error := fmt.Errorf("couldn't get payload from authtoken")
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	arg := &pb.UserId{
		Id: payload.UserId,
	}
	user, err := c.GetUser(context, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

type userIdParam struct {
	ID int64 `uri:"id"  binding:"required,min=1"`
}

type updateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	var params userIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	arg := &pb.UpdateUserRequest{
		Id:   params.ID,
		Name: req.Name,
	}

	user, err := c.UpdateUser(context, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type addPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (server *Server) addPhone(ctx *gin.Context) {
	var req addPhoneRequest
	var params userIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	arg := &pb.AddPhoneRequest{
		Id:    params.ID,
		Phone: req.Phone,
	}

	_, err = c.AddPhone(context, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "Phone has been added")
}

type deleteUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		error := fmt.Errorf("failed to connect to UsersService: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	arg := &pb.UserId{
		Id: req.ID,
	}

	_, err = c.DeleteUser(context, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "User was deleted")
}

func (server *Server) logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "logged out")
}
