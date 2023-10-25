package main

import (
	"log"

	"github.com/zura-t/go_delivery_system/config"
	"github.com/zura-t/go_delivery_system/internal/app"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config file:", err)
	}

	app.Run(config)
}
