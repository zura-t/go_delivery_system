package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zura-t/go_delivery_system/cmd/api"
	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/rmq"
)

func main() {
	rabbitConn, err := connectRabbitmq()
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	channel, emitter, err := setupRabbitmq(rabbitConn)
	defer channel.Close()

	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config file:", err)
	}

	runGinServer(config, rabbitConn, emitter)
}

func runGinServer(config internal.Config, rabbitConn *amqp.Connection, emitter *rmq.Emitter) {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	server, err := api.NewServer(config, rabbitConn, emitter)
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

func connectRabbitmq() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			fmt.Println("rabbitmq not yet ready")
			counts++
		} else {
			log.Println("Connected to rabbitmq.")

			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("back off")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func setupRabbitmq(rabbitConn *amqp.Connection) (*amqp.Channel, *rmq.Emitter, error) {
	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("can't create rabbitmq emitter", err)
		return nil, nil, err
	}

	emitter, err := rmq.NewEmitter(rabbitConn, channel)
	if err != nil {
		return nil, nil, err
	}

	return channel, &emitter, nil
}
