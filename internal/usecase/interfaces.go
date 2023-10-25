package usecase

import (
	"github.com/zura-t/go_delivery_system/internal/entity"
)

type User interface {
	CreateUser(req entity.UserRegister) (entity.User, int, error)
	LoginUser(req entity.UserLogin) (entity.UserLoginResponse, int, error)
	GetMyProfile(id int64) (entity.User, int, error)
	UpdateUser(id int64, req entity.UserUpdate) (entity.User, int, error)
	AddPhone(id int64, req entity.UserAddPhone) (string, int, error)
	DeleteUser(id int64) (string, int, error)
}
