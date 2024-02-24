package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/entity"
	"github.com/zura-t/go_delivery_system/internal/usecase/httpclient"
	"github.com/zura-t/go_delivery_system/pkg/httpserver"
)

type ShopWebAPI struct {
	client *http.Client
	config *config.Config
}

func NewShopWebAPI(config *config.Config) *ShopWebAPI {
	return &ShopWebAPI{
		client: &http.Client{},
		config: config,
	}
}

func (webapi *ShopWebAPI) CreateShop(req *entity.CreateShop) (*entity.Shop, int, error) {
	url := fmt.Sprintf("%s/shops", webapi.config.ShopsServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPost, url)
	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.Shop{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.Shop{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var shop *entity.Shop
	newShop, err := io.ReadAll(res.Body)

	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newShop, &shop)
	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	return shop, http.StatusOK, nil
}

func (webapi *ShopWebAPI) GetShopInfo(id int64) (*entity.Shop, int, error) {
	url := fmt.Sprintf("%s/shops/%d", webapi.config.ShopsServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodGet, url)
	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return &entity.Shop{}, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return &entity.Shop{}, res.StatusCode, err
	}
	defer res.Body.Close()

	var shop *entity.Shop
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &shop)
	if err != nil {
		return &entity.Shop{}, http.StatusInternalServerError, err
	}
	return shop, http.StatusOK, nil
}

func (webapi *ShopWebAPI) GetShops(limit int32, offset int32) ([]*entity.Shop, int, error) {
	url := fmt.Sprintf("%s/shops", webapi.config.ShopsServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodGet, url)
	query := httpRequest.URL.Query()
	query.Add("limit", strconv.Itoa(int(limit)))
	query.Add("offset", strconv.Itoa(int(offset)))
	httpRequest.URL.RawQuery = query.Encode()

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var shops []entity.Shop
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &shops)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	response := make([]*entity.Shop, len(shops))
	for i := 0; i < len(shops); i++ {
		response[i] = &shops[i]
	}
	return response, http.StatusOK, nil
}

func (webapi *ShopWebAPI) GetShopsAdmin(user_id int64) ([]entity.Shop, int, error) {
	url := fmt.Sprintf("%s/shops/admin", webapi.config.ShopsServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodGet, url)
	query := httpRequest.URL.Query()
	query.Add("user_id", strconv.Itoa(int(user_id)))
	httpRequest.URL.RawQuery = query.Encode()

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var shops []entity.Shop
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &shops)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return shops, http.StatusOK, nil
}

func (webapi *ShopWebAPI) UpdateShop(id int64, req *entity.UpdateShopInfo) (*entity.Shop, int, error) {
	url := fmt.Sprintf("%s/shops/%d", webapi.config.ShopsServiceAddress, id)

	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPatch, url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var shop *entity.Shop
	newShop, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newShop, &shop)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return shop, http.StatusOK, nil
}

func (webapi *ShopWebAPI) CreateMenu(req []*entity.CreateMenuItem) ([]*entity.MenuItem, int, error) {
	url := fmt.Sprintf("%s/shops/menu_items", webapi.config.ShopsServiceAddress)

	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPost, url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var menuItems []*entity.MenuItem
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &menuItems)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return menuItems, http.StatusOK, nil
}

func (webapi *ShopWebAPI) GetMenu(shopId int64) ([]*entity.MenuItem, int, error) {
	url := fmt.Sprintf("%s/shops/menu_items", webapi.config.ShopsServiceAddress)
	httpRequest, err := httpclient.NewHttpRequest(shopId, http.MethodGet, url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var menuItems []*entity.MenuItem
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &menuItems)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return menuItems, http.StatusOK, nil
}

func (webapi *ShopWebAPI) UpdateMenuItem(id int64, req *entity.UpdateMenuItem) (*entity.MenuItem, int, error) {
	url := fmt.Sprintf("%s/shops/menu_items/%d", webapi.config.ShopsServiceAddress, id)

	httpRequest, err := httpclient.NewHttpRequest(req, http.MethodPatch, url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var menuItem *entity.MenuItem
	newMenuItem, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(newMenuItem, &menuItem)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return menuItem, http.StatusOK, nil
}

func (webapi *ShopWebAPI) GetMenuItem(id int64) (*entity.MenuItem, int, error) {
	url := fmt.Sprintf("%s/shops/menu_items/%d", webapi.config.ShopsServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodGet, url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := webapi.client.Do(httpRequest)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if res.StatusCode != 200 {
		errorMessage, err := httpserver.HttpErrorResponse(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = fmt.Errorf("Error: %s", errorMessage)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	var menuItem *entity.MenuItem
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(resp, &menuItem)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return menuItem, http.StatusOK, nil
}

func (webapi *ShopWebAPI) DeleteShop(id int64, user_id int64) (string, int, error) {
	url := fmt.Sprintf("%s/shops/%d", webapi.config.ShopsServiceAddress, id)
	httpRequest, err := httpclient.NewHttpRequest(nil, http.MethodDelete, url)
	query := httpRequest.URL.Query()
	query.Add("user_id", strconv.Itoa(int(user_id)))
	httpRequest.URL.RawQuery = query.Encode()

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

func (webapi *ShopWebAPI) DeleteMenuItem(id int64) (string, int, error) {
	url := fmt.Sprintf("%s/shops/menu_items/%d", webapi.config.ShopsServiceAddress, id)
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
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}
