package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/internal/entity"
	"github.com/zura-t/go_delivery_system/internal/usecase"
	"github.com/zura-t/go_delivery_system/pkg/logger"
)

type shopRoutes struct {
	shopUsecase usecase.Shop
	logger      logger.Interface
}

func (server *Server) newShopRoutes(handler *gin.RouterGroup, shopUsecase usecase.Shop, logger logger.Interface) {
	routes := &shopRoutes{shopUsecase, logger}

	handler.Group("/").Use(authMiddleware(server.tokenMaker))
	
	shopRoutes := handler.Group("/shops")
	menuItemRoutes := shopRoutes.Group("/menu_items")

	shopRoutes.POST("/", routes.createShop).Use(rolesMiddleware())
	shopRoutes.GET("/:id", routes.getShop)
	shopRoutes.GET("/", routes.getShops)
	shopRoutes.GET("/admin", routes.getShopsAdmin).Use(rolesMiddleware())
	shopRoutes.PATCH("/:id", routes.updateShop).Use(rolesMiddleware())
	shopRoutes.DELETE("/:id", routes.deleteShop).Use(rolesMiddleware())

	menuItemRoutes.POST("/", routes.createMenuItems).Use(rolesMiddleware())
	menuItemRoutes.GET("/:shop_id", routes.getMenuItems)
	menuItemRoutes.PATCH("/:id", routes.updateMenuItem).Use(rolesMiddleware())
	menuItemRoutes.GET("/:id", routes.getMenuItem)
	menuItemRoutes.DELETE("/:id", routes.deleteMenuItem).Use(rolesMiddleware())
}

type CreateShopRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	OpenTime    time.Time `json:"open_time" binding:"required"`
	CloseTime   time.Time `json:"close_time" binding:"required"`
	IsClosed    bool      `json:"is_closed" binding:"required"`
}

func (r *shopRoutes) createShop(ctx *gin.Context) {
	var req CreateShopRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - createUser")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.shopUsecase.CreateShop(&entity.CreateShop{
		Name:        req.Name,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		IsClosed:    req.IsClosed,
	})
	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createShop")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type IdParam struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (r *shopRoutes) getShop(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getShop")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	shop, st, err := r.shopUsecase.GetShopInfo(req.Id)

	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getMyProfile")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shop)
}

func (r *shopRoutes) getShops(ctx *gin.Context) {
	shops, st, err := r.shopUsecase.GetShops()

	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getShops")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shops)
}

type GetShopsAdmin struct {
	UserId int64 `uri:"id" binding:"required,min=1"`
}

func (r *shopRoutes) getShopsAdmin(ctx *gin.Context) {
	payload := getJWTPayload(ctx)

	shops, st, err := r.shopUsecase.GetShopsAdmin(payload.UserId)

	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getShopsAdmin")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shops)
}

type UpdateShopRequest struct {
	Id          int64     `uri:"id" binding:"required,min=1"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	IsClosed    bool      `json:"is_closed"`
}

func (r *shopRoutes) updateShop(ctx *gin.Context) {
	var req UpdateShopRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateShop")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.shopUsecase.UpdateShop(req.Id, &entity.UpdateShopInfo{
		Name:        req.Name,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		IsClosed:    req.IsClosed,
	})

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateShop")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type CreateMenuItemsRequest struct {
	MenuItems []struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Photo       string `json:"photo"`
		Price       int32  `json:"price" binding:"required,min=1"`
		ShopID      int64  `json:"shop_id" binding:"required,min=1"`
	} `json:"menu_items"`
}

func (r *shopRoutes) createMenuItems(ctx *gin.Context) {
	var req CreateMenuItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var request []*entity.CreateMenuItem
	for i := 0; i < len(req.MenuItems); i++ {
		request[i] = &entity.CreateMenuItem{
			Name:        req.MenuItems[i].Name,
			Description: req.MenuItems[i].Description,
			Photo:       req.MenuItems[i].Photo,
			Price:       req.MenuItems[i].Price,
			ShopID:      req.MenuItems[i].ShopID,
		}
	}

	menuItems, st, err := r.shopUsecase.CreateMenu(request)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuItems)
}

type GetMenuRequest struct {
	ShopId int64 `uri:"shop_id" binding:"required,min=1"`
}

func (r *shopRoutes) getMenuItems(ctx *gin.Context) {
	var req GetMenuRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	menuItems, st, err := r.shopUsecase.GetMenu(req.ShopId)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuItems)
}

type UpdateMenuItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price" binding:"min=1"`
}

func (r *shopRoutes) updateMenuItem(ctx *gin.Context) {
	var req UpdateMenuItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var params IdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	menuItems, st, err := r.shopUsecase.UpdateMenuItem(params.Id, &entity.UpdateMenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuItems)
}

func (r *shopRoutes) getMenuItem(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getMenuItem")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	menuItem, st, err := r.shopUsecase.GetMenuItem(req.Id)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getMenuItem")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuItem)
}

func (r *shopRoutes) deleteShop(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, st, err := r.shopUsecase.DeleteShop(req.Id)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (r *shopRoutes) deleteMenuItem(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, st, err := r.shopUsecase.DeleteMenuItem(req.Id)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}
