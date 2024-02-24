package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/entity"
	"github.com/zura-t/go_delivery_system/internal/usecase/httpclient"
	"github.com/zura-t/go_delivery_system/pkg/httpserver"
)

type UserWebAPI struct {
	client *http.Client
	config *config.Config
}

func NewUserWebAPI(config *config.Config) *UserWebAPI {
	return &UserWebAPI{
		client: &http.Client{},
		config: config,
	}
}

func (webapi *UserWebAPI) CreateUser(req *entity.UserRegister) (*entity.User, int, error) {
	url := fmt.Sprintf("%s/users", webapi.config.UsersServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPost, url)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}

	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	newUser, err := io.ReadAll(res.Body)

	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (webapi *UserWebAPI) LoginUser(req *entity.UserLogin) (*entity.UserLoginResponse, int, error) {
	url := fmt.Sprintf("%s/login", webapi.config.UsersServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPost, url)
	if err != nil {
		return &entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.UserLoginResponse{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.UserLoginResponse{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.UserLoginResponse
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return &entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		return &entity.UserLoginResponse{}, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (webapi *UserWebAPI) GetMyProfile(id int64) (*entity.User, int, error) {
	url := fmt.Sprintf("%s/users/my_profile/%d", webapi.config.UsersServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodGet, url)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (webapi *UserWebAPI) AddAdminRole(id int64) (string, int, error) {
	url := fmt.Sprintf("%s/users/admin/%d", webapi.config.UsersServiceAddress, id)

	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodPatch, url)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

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
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (webapi *UserWebAPI) UpdateUser(id int64, req *entity.UserUpdate) (*entity.User, int, error) {
	url := fmt.Sprintf("%s/users/%d", webapi.config.UsersServiceAddress, id)

	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPatch, url)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.User{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.User{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var user entity.User
	newUser, err := io.ReadAll(res.Body)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newUser, &user)
	if err != nil {
		return &entity.User{}, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

func (webapi *UserWebAPI) AddPhone(id int64, req *entity.UserAddPhone) (string, int, error) {
	url := fmt.Sprintf("%s/users/phone_number/%d", webapi.config.UsersServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPatch, url)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

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
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (webapi *UserWebAPI) DeleteUser(id int64) (string, int, error) {
	url := fmt.Sprintf("%s/users/%d", webapi.config.UsersServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodDelete, url)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

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
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}
