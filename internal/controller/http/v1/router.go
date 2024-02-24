package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/internal/usecase"
	"github.com/zura-t/go_delivery_system/pkg/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zura-t/go_delivery_system/docs"
)

func (server *Server) NewRouter(handler *gin.Engine, logger logger.Interface, userUsecase usecase.User, shopsUsecase usecase.Shop) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	{
		server.newUserRoutes(handler, userUsecase, logger)
		server.newShopRoutes(handler, shopsUsecase, logger)
	}
}
