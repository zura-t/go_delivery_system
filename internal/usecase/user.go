package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/entity"
	"github.com/zura-t/go_delivery_system/pkg/httpserver"
)

type UserUseCase struct {
	config *config.Config
}

func New(config *config.Config) *UserUseCase {
	return &UserUseCase{
		config: config,
	}
}

func (uc *UserUseCase) CreateUser(req entity.UserRegister) (entity.User, int, error) {
	url := fmt.Sprintf("%s/users", uc.config.UsersServiceAddress)
	res, err := httpserver.SendHttpRequest(req, http.MethodPost, url)

	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	newUser, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (uc *UserUseCase) LoginUser(req entity.UserLogin) (entity.UserLoginResponse, int, error) {
	url := fmt.Sprintf("%s/users/login", uc.config.UsersServiceAddress)
	res, err := httpserver.SendHttpRequest(req, http.MethodPost, url)

	if err != nil {
		return entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return entity.UserLoginResponse{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return entity.UserLoginResponse{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.UserLoginResponse
	newUser, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (uc *UserUseCase) GetMyProfile(id int64) (entity.User, int, error) {
	url := fmt.Sprintf("%s/users/%d", uc.config.UsersServiceAddress, id)
	res, err := httpserver.SendHttpRequest(nil, http.MethodGet, url)

	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	newUser, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (uc *UserUseCase) UpdateUser(id int64, req entity.UserUpdate) (entity.User, int, error) {
	url := fmt.Sprintf("%s/users/%d", uc.config.UsersServiceAddress, id)

	res, err := httpserver.SendHttpRequest(req, http.MethodPatch, url)

	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	newUser, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (uc *UserUseCase) AddPhone(id int64, req entity.UserAddPhone) (string, int, error) {
	url := fmt.Sprintf("%s/users/phone_number/%d", uc.config.UsersServiceAddress, id)
	res, err := httpserver.SendHttpRequest(req, http.MethodPatch, url)

	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return "", res.StatusCode, err
	}
	defer res.Body.Close()

	var resp string
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (uc *UserUseCase) DeleteUser(id int64) (string, int, error) {
	url := fmt.Sprintf("%s/users/%d", uc.config.UsersServiceAddress, id)
	res, err := httpserver.SendHttpRequest(nil, http.MethodDelete, url)

	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return "", res.StatusCode, err
	}
	defer res.Body.Close()

	var resp string
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}
