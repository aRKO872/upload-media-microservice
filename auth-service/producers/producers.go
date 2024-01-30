package producers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upload-media-auth/types"
)

func failOnError(err error, msg string) {
  if err != nil {
    fmt.Printf("%s: %s", msg, err)
  }
}

func PublishToQueue(queueName string, msg types.ProducedOrConsumedMessage) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	messageBody, err := json.Marshal(msg)
	if err != nil {
		failOnError(err, "Error serializing object")
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "application/json",
			Body:        messageBody,
		},
	)

	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s\n", msg)
}