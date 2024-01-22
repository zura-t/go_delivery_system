package v1

import (
	"fmt"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/pkg/logger"
	"github.com/zura-t/go_delivery_system/token"
)

type Server struct {
	config     *config.Config
	tokenMaker token.Maker
	l          *logger.Logger
}

func New(cfg *config.Config, l *logger.Logger) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}
	return &Server{
		config:     cfg,
		tokenMaker: tokenMaker,
		l:          l,
	}, nil
}
