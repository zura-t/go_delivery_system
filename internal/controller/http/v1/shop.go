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

func (server *Server) newShopRoutes(handler *gin.Engine, shopUsecase usecase.Shop, logger logger.Interface) {
	routes := &shopRoutes{shopUsecase, logger}

	handler.Group("/").Use(authMiddleware(server.tokenMaker))

	shopRoutes := handler.Group("/shops")
	menuItemRoutes := shopRoutes.Group("/menu_items")

	shopRoutes.POST("/", routes.createShop).Use(server.rolesMiddleware())
	shopRoutes.GET("/:id", routes.getShop)
	shopRoutes.GET("/", routes.getShops)
	shopRoutes.GET("/admin", routes.getShopsAdmin).Use(server.rolesMiddleware())
	shopRoutes.PATCH("/:id", routes.updateShop).Use(server.rolesMiddleware())
	shopRoutes.DELETE("/:id", routes.deleteShop).Use(server.rolesMiddleware())

	menuItemRoutes.POST("/", routes.createMenuItems).Use(server.rolesMiddleware())
	menuItemRoutes.GET("/list/:id", routes.getMenuItems)
	menuItemRoutes.PATCH("/:id", routes.updateMenuItem).Use(server.rolesMiddleware())
	menuItemRoutes.GET("/:id", routes.getMenuItem)
	menuItemRoutes.DELETE("/:id", routes.deleteMenuItem).Use(server.rolesMiddleware())
}

type CreateShopRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" example:""`
	OpenTime    time.Time `json:"open_time" binding:"required" example:""`
	CloseTime   time.Time `json:"close_time" binding:"required" example:""`
	IsClosed    bool      `json:"is_closed"`
}

// @Summary     Create Shop
// @Description Create new Shop
// @ID          create-shop
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       request body CreateShopRequest true "CreateShop"
// @Success     200 {object} entity.Shop
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/ [post]
func (r *shopRoutes) createShop(ctx *gin.Context) {
	var req CreateShopRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - user routes - createUser")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := getJWTPayload(ctx)

	shop, st, err := r.shopUsecase.CreateShop(&entity.CreateShop{
		Name:        req.Name,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		UserId:      payload.UserId,
		IsClosed:    req.IsClosed,
	})
	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createShop")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shop)
}

type IdParam struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// @Summary     Get Shop
// @Description Get Shop info
// @ID          getShop
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       request path IdParam true "getShop"
// @Success     200 {object} entity.Shop
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/{id} [get]
func (r *shopRoutes) getShop(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - getShop")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	shop, st, err := r.shopUsecase.GetShop(req.Id)

	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getMyProfile")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shop)
}

type GetShopsRequest struct {
	Limit  int32 `form:"limit,default=20" binding:"min=1"`
	Offset int32 `form:"offset,default=0"`
}

// @Summary     GetShops
// @Description getShops
// @ID          getShops
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.Shop
// @Param       limit query string false "rows to return"
// @Param       offset query string  false  "rows to skip"
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/ [get]
func (r *shopRoutes) getShops(ctx *gin.Context) {
	var req GetShopsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	shops, st, err := r.shopUsecase.GetShops(req.Limit, req.Offset)

	if err != nil {
		r.logger.Error(err, "http - v1 - user routes - getShops")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, shops)
}

// @Summary     GetShopsAdmin
// @Description get shops where you're admin
// @ID          getShopsAdmin
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.Shop
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/admin [get]
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
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	IsClosed    bool      `json:"is_closed"`
}

// @Summary     Update Shop
// @Description Update Shop
// @ID          update-shop
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       request body UpdateShopRequest true "updateShop"
// @Param       id path IdParam true "id"
// @Success     200 {object} entity.Shop
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/{id} [patch]
func (r *shopRoutes) updateShop(ctx *gin.Context) {
	var req UpdateShopRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateShop")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var param IdParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateShop")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := getJWTPayload(ctx)

	data, st, err := r.shopUsecase.UpdateShop(param.Id, &entity.UpdateShopInfo{
		Name:        req.Name,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		IsClosed:    req.IsClosed,
		UserId:      payload.UserId,
	})

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateShop")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type MenuItem struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Price       int32  `json:"price" binding:"required,min=1"`
}

type CreateMenuItemsRequest struct {
	MenuItems []MenuItem `json:"menu_items" binding:"required,min=1"`
	ShopId    int64      `json:"shop_id" binding:"required,min=1"`
}

// @Summary     Create MenuItems
// @Description Create MenuItems
// @ID          create-menuitems
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       request body CreateMenuItemsRequest true "create menuitems"
// @Success     200 {object} []entity.MenuItem
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/menu_items [post]
func (r *shopRoutes) createMenuItems(ctx *gin.Context) {
	var req CreateMenuItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := getJWTPayload(ctx)

	menuItems := make([]entity.MenuItem, len(req.MenuItems))
	for i := 0; i < len(req.MenuItems); i++ {
		menuItems[i] = entity.MenuItem{
			Name:        req.MenuItems[i].Name,
			Description: req.MenuItems[i].Description,
			Price:       req.MenuItems[i].Price,
		}
	}

	menuCreated, st, err := r.shopUsecase.CreateMenu(&entity.CreateMenuItem{
		MenuItems: menuItems,
		ShopId:    req.ShopId,
		UserId:    payload.UserId,
	})

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - createMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuCreated)
}

type GetMenuRequest struct {
	ShopId int64 `uri:"id" binding:"required,min=1"`
}

// @Summary     GetMenuItems
// @Description getMenuItems
// @ID          getMenuItems
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       id path IdParam true "id"
// @Success     200 {object} []entity.MenuItem
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/menu_items/list/{id} [get]
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

// @Summary     updateMenuItem
// @Description updateMenuItem
// @ID          updateMenuItem
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       request path IdParam true "id"
// @Param       request body UpdateMenuItemRequest true "updateMenuItem"
// @Success     200 {object} entity.MenuItem
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/menu_items/{id} [patch]
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

	payload := getJWTPayload(ctx)

	menuItems, st, err := r.shopUsecase.UpdateMenuItem(params.Id, &entity.UpdateMenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UserId:      payload.UserId,
	})

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - updateMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, menuItems)
}

// @Summary     getMenuItem
// @Description getMenuItem
// @ID          getMenuItem
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       id path IdParam true "id"
// @Success     200 {object} entity.MenuItem
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/menu_items/{id} [get]
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

type UserIdQuery struct {
	UserId int64 `form:"user_id" binding:"required,min=1"`
}

// @Summary     Delete Shop
// @Description Delete Shop
// @ID          delete-shop
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       id path IdParam true "id"
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/{id} [delete]
func (r *shopRoutes) deleteShop(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := getJWTPayload(ctx)

	res, st, err := r.shopUsecase.DeleteShop(req.Id, payload.UserId)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary     DeleteMenuItem
// @Description DeleteMenuItem
// @ID          deleteMenuItem
// @Tags  	    shops
// @Accept      json
// @Produce     json
// @Param       id path IdParam true "id"
// @Success     200 {object} string
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Security 		BearerAuth
// @Router      /shops/menu_items/{id} [delete]
func (r *shopRoutes) deleteMenuItem(ctx *gin.Context) {
	var req IdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := getJWTPayload(ctx)

	res, st, err := r.shopUsecase.DeleteMenuItem(req.Id, payload.UserId)

	if err != nil {
		r.logger.Error(err, "http - v1 - shop routes - deleteMenuItems")
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}
