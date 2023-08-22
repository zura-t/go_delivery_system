package gapi

import (
	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/pb"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	config internal.Config
}

func NewServer(config internal.Config) (*Server, error) {
	server := &Server{config: config}
	return server, nil
}
