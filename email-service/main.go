package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	emailcontroller "github.com/upload-media-email/controllers"
	"github.com/upload-media-email/types"
)

var emailConn *amqp.Connection
var emailChannel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	e := echo.New()

	// Setup RabbitMQ consumer
	var err error
	emailConn, err = amqp.Dial("amqp://guest:guest@rabbitmq")
	failOnError(err, "Failed to connect to RabbitMQ")

	// defer emailConn.Close()

	emailChannel, err = emailConn.Channel()
	failOnError(err, "Failed to open a channel")

	// defer emailChannel.Close()

	// Declare the queue
	q, err := emailChannel.QueueDeclare(
		"email-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	msgs, err := emailChannel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	// Goroutine to process messages
	go func() {
		for d := range msgs {
			var emailRequest types.ProducedOrConsumedMessage
			err := json.Unmarshal(d.Body, &emailRequest)
			if err != nil {
				fmt.Printf("Failed to unmarshal message body: %v", err)
			}

			emailcontroller.SendEmail (emailRequest.Email, emailRequest.PictureURL)
		}
	}()

	log.Println("RabbitMQ consumer setup completed")

	// Start the Echo server
	e.Start(":8087")

	// Graceful shutdown (if needed)
	gracefulShutdown(e)
}

func gracefulShutdown(e *echo.Echo) {
	// Perform cleanup tasks on application shutdown
	// For example, close RabbitMQ connection
	if emailConn != nil {
		emailChannel.Close()
		emailConn.Close()
	}
	log.Println("Server stopped gracefully")
	time.Sleep(1 * time.Second)
}
