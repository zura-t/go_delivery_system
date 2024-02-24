package usecase

import (
	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/entity"
)

type UserUseCase struct {
	config *config.Config
	webapi UserWebAPI
}

func NewUserUseCase(config *config.Config, webapi UserWebAPI) *UserUseCase {
	return &UserUseCase{
		config: config,
		webapi: webapi,
	}
}

func (uc *UserUseCase) CreateUser(req *entity.UserRegister) (*entity.User, int, error) {
	return uc.webapi.CreateUser(req)
}

func (uc *UserUseCase) LoginUser(req *entity.UserLogin) (*entity.UserLoginResponse, int, error) {
	return uc.webapi.LoginUser(req)
}

func (uc *UserUseCase) GetMyProfile(id int64) (*entity.User, int, error) {
	return uc.webapi.GetMyProfile(id)
}

func (uc *UserUseCase) AddAdminRole(id int64) (string, int, error) {
	return uc.webapi.AddAdminRole(id)
}

func (uc *UserUseCase) UpdateUser(id int64, req *entity.UserUpdate) (*entity.User, int, error) {
	return uc.webapi.UpdateUser(id, req)
}

func (uc *UserUseCase) AddPhone(id int64, req *entity.UserAddPhone) (string, int, error) {
	return uc.webapi.AddPhone(id, req)
}

func (uc *UserUseCase) DeleteUser(id int64) (string, int, error) {
	return uc.webapi.DeleteUser(id)
}

func IsAdmin(id int64) (bool, error) {
	return true, nil
}