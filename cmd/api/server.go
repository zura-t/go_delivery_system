package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/token"
)

type Server struct {
	router     *gin.Engine
	config     internal.Config
	tokenMaker token.Maker
}

func NewServer(config internal.Config) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/logout", server.logout)
	router.POST("/renew_token", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/my_profile", server.getMyProfile)
	authRoutes.PATCH("/users/:id", server.updateUser)
	authRoutes.PATCH("/users/phone_number/:id", server.addPhone)
	authRoutes.DELETE("/users/:id", server.deleteUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
