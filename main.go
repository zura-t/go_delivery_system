package main

import (
	"log"

	"github.com/gin-gonic/gin"
	// "github.com/zura-t/go_delivery_system/cmd/api"
	"github.com/zura-t/go_delivery_system/cmd/api"
	"github.com/zura-t/go_delivery_system/internal"
)

func main() {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config file:", err)
	}

	runGinServer(config)
}

func runGinServer(config internal.Config) {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatalf("can't create server: %s", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatalf("can't start server: %s", err)
	}

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatalf("cannot create statik fs: %s", err)
	// }
	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

}
