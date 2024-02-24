package rmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func SendRPC[Q, S any](routingKey string, r Q) (rsp *S, err error) {
	rabbitConn, err := amqp.Dial("amqp://admin:admin@localhost:5672")
	if err != nil {
		log.Println(err)
	}
	defer rabbitConn.Close()

	ch, err := rabbitConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue
		"",                      // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return nil, err
	}

	corrId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	req, _ := json.Marshal(r)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId.String(),
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          []byte(req),
		},
	)
	if err != nil {
		return nil, err
	}

	for d := range msgs {
		if corrId.String() == d.CorrelationId {
			err = json.Unmarshal(d.Body, &rsp)
			log.Println(rsp)
			break
		}
		break
	}

	return
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type CreateUserPayload struct {
	Name string
	Data CreateUserRequest
}

func NewEmitter(conn *amqp.Connection, channel *amqp.Channel) (Emitter, error) {
	emitter := Emitter{
		conn:    conn,
		channel: channel,
	}

	return emitter, nil
}
