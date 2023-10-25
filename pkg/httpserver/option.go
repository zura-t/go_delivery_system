package httpserver

import (
	"net"
	"time"
)

type Option func(*HttpServer)

func Port(port string) Option {
	return func(s *HttpServer) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.shutdownTimeout = timeout
	}
}