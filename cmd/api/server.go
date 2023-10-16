package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/rmq"
	"github.com/zura-t/go_delivery_system/token"
)

type Server struct {
	router     *gin.Engine
	config     internal.Config
	tokenMaker token.Maker
	rabbit     *amqp.Connection
	emitter    *rmq.Emitter
}

func NewServer(config internal.Config, rabbitConn *amqp.Connection, emitter *rmq.Emitter) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		rabbit:     rabbitConn,
		emitter:    emitter,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/users", server.CreateUser)
	router.POST("/login", server.LoginUser)
	router.POST("/logout", server.Logout)
	router.POST("/renew_token", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/my_profile", server.GetMyProfile)
	authRoutes.PATCH("/users/:id", server.UpdateUser)
	authRoutes.PATCH("/users/phone_number/:id", server.AddPhone)
	authRoutes.DELETE("/users/:id", server.DeleteUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func httpErrorResponse(body io.ReadCloser) (gin.H, error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var errorMessage gin.H
	err = json.Unmarshal(content, &errorMessage)
	if err != nil {
		return nil, err
	}

	return errorMessage, nil
}
