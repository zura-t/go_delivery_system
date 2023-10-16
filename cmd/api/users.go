package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/token"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type UserResponse struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := fmt.Sprintf("%s/users", server.config.UsersServiceAddress)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}
	defer res.Body.Close()

	var user UserResponse
	newUser, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := fmt.Sprintf("%s/users/login", server.config.UsersServiceAddress)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}
	defer res.Body.Close()

	var user LoginUserResponse
	response, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(response, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.SetCookie("refresh_token", user.RefreshToken, int(time.Until(user.RefreshTokenExpiresAt).Seconds()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, user)
}

func (server *Server) GetMyProfile(ctx *gin.Context) {
	var payload token.Payload
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

	url := fmt.Sprintf("%s/users/%d", server.config.UsersServiceAddress, payload.UserId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}
	defer res.Body.Close()

	var user UserResponse
	response, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(response, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type UserIdParam struct {
	Id int64 `uri:"id"  binding:"required,min=1"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	var params UserIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newUser := &UpdateUserRequest{
		Name: req.Name,
	}

	arg, err := json.Marshal(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := fmt.Sprintf("%s/users/%d", server.config.UsersServiceAddress, params.Id)
	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(arg))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}
	defer res.Body.Close()

	var user UserResponse
	response, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(response, &user)

	ctx.JSON(http.StatusOK, user)
}

type AddPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (server *Server) AddPhone(ctx *gin.Context) {
	var req AddPhoneRequest
	var params UserIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userData := &AddPhoneRequest{
		Phone: req.Phone,
	}

	arg, err := json.Marshal(userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := fmt.Sprintf("%s/users/phone_number/%d", server.config.UsersServiceAddress, params.Id)
	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(arg))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}

	ctx.JSON(http.StatusOK, "Phone has been added")
}

type DeleteUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) DeleteUser(ctx *gin.Context) {
	var req DeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	url := fmt.Sprintf("%s/users/%d", server.config.UsersServiceAddress, req.Id)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpErrorResponse(res.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(res.StatusCode, errorMessage)
		return
	}

	ctx.JSON(http.StatusOK, "User was deleted")
}

func (server *Server) Logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "logged out")
}
