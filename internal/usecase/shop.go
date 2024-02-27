package usecase

import (
	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/entity"
)

type ShopUseCase struct {
	config *config.Config
	webapi ShopWebAPI
}

func NewShopUseCase(config *config.Config, webapi ShopWebAPI) *ShopUseCase {
	return &ShopUseCase{
		config: config,
		webapi: webapi,
	}
}

func (uc *ShopUseCase) CreateShop(req *entity.CreateShop) (*entity.Shop, int, error) {
	return uc.webapi.CreateShop(req)
}

func (uc *ShopUseCase) GetShop(id int64) (*entity.Shop, int, error) {
	return uc.webapi.GetShopInfo(id)
}

func (uc *ShopUseCase) GetShops(limit int32, offset int32) ([]*entity.Shop, int, error) {
	return uc.webapi.GetShops(limit, offset)
}

func (uc *ShopUseCase) GetShopsAdmin(user_id int64) ([]entity.Shop, int, error) {
	return uc.webapi.GetShopsAdmin(user_id)
}

func (uc *ShopUseCase) UpdateShop(id int64, req *entity.UpdateShopInfo) (*entity.Shop, int, error) {
	return uc.webapi.UpdateShop(id, req)
}

func (uc *ShopUseCase) CreateMenu(req *entity.CreateMenuItem) ([]*entity.GetMenuItem, int, error) {
	return uc.webapi.CreateMenu(req)
}

func (uc *ShopUseCase) GetMenu(shopId int64) ([]*entity.GetMenuItem, int, error) {
	return uc.webapi.GetMenu(shopId)
}

func (uc *ShopUseCase) UpdateMenuItem(id int64, req *entity.UpdateMenuItem) (*entity.GetMenuItem, int, error) {
	return uc.webapi.UpdateMenuItem(id, req)
}

func (uc *ShopUseCase) GetMenuItem(id int64) (*entity.GetMenuItem, int, error) {
	return uc.webapi.GetMenuItem(id)
}

func (uc *ShopUseCase) DeleteShop(id int64, user_id int64) (string, int, error) {
	return uc.webapi.DeleteShop(id, user_id)
}

func (uc *ShopUseCase) DeleteMenuItem(id int64) (string, int, error) {
	return uc.webapi.DeleteMenuItem(id)
}
