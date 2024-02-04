package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system/config"
	v1 "github.com/zura-t/go_delivery_system/internal/controller/http/v1"
	"github.com/zura-t/go_delivery_system/internal/usecase"
	"github.com/zura-t/go_delivery_system/internal/usecase/webapi"
	"github.com/zura-t/go_delivery_system/pkg/httpserver"
	"github.com/zura-t/go_delivery_system/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		l.Fatal("Connection failed")
		os.Exit(1)
	}
	_ = esClient
	webapi := webapi.NewUserWebAPI(cfg)
	usersUseCase := usecase.NewUserUseCase(cfg, webapi)

	// rmqRouter := amqprpc.NewRouter(translationUseCase)

	// rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	// }

	runGinServer(l, cfg, usersUseCase)
}

func runGinServer(l *logger.Logger, cfg *config.Config, usersUseCase *usecase.UserUseCase) {
	handler := gin.New()
	server, err := v1.New(cfg, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - runGinServer: %w", err))
		os.Exit(1)
	}
	server.NewRouter(handler, l, usersUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HttpPort))
	log.Print(fmt.Sprintf("server started on port %s", cfg.HttpPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print(s.String())
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Fatal(err)
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err := <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Fatal(err)
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	// err = rmqServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	// }

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatalf("cannot create statik fs: %s", err)
	// }
	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

}
