package v1

import (
	"fmt"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/token"
)

type Server struct {
	config     *config.Config
	tokenMaker token.Maker
}

func New(cfg *config.Config) (*Server, error){
	tokenMaker, err := token.NewJwtMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}
	return &Server{
		config:     cfg,
		tokenMaker: tokenMaker,
	}, nil
}
