package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/internal/usecase"
	"github.com/zura-t/go_delivery_system/pkg/logger"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger spec:
// @title       Microservice api-gateway
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func (server *Server) NewRouter(handler *gin.Engine, logger logger.Interface, userUsecase usecase.User) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	// handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.Group("/v1")
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	{
		server.newUserRoutes(h, userUsecase, logger)
	}
}
