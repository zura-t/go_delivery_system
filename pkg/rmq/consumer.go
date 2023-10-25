package rmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	consumer.channel = ch
	return nil
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	q, err := consumer.channel.QueueDeclare(
		"queue.", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return err
	}

	for _, s := range topics {
		consumer.channel.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := consumer.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			// var payload Payload
			// _ = json.Unmarshal(d.Body, &payload)
			log.Printf("Received a message: %s", d.Body)
			log.Printf("Done")

			// go handlePayload(payload)
		}
	}()

	fmt.Printf("Wait for message on [Exchange, Queue ] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}
