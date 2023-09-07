package gapi

import (
	"fmt"

	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/pb"
	"github.com/zura-t/go_delivery_system/token"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
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
	return server, nil
}
