package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/faraji-fuji/learning-rabbitmq/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {
	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Printf("Something went wrong: %s", err)
		return
	}
	logger.Printf("Connected to RabbitMQ")
	defer conn.Close()

	// create channel
	ch, err := pubsub.InitPubSub(conn, "test-exchange", "test-queue")
	if err != nil {
		logger.Printf("Something went wrong: %s", err)
		return
	}
	logger.Printf("Created channel")
	defer ch.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish with context
	body := "Hello..."
	err = ch.PublishWithContext(ctx,
		"test-exchange",
		"test-queue",
		false,
		false,
		amqp.Publishing{
			Headers:     amqp.Table{"x-delay": 5000},
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		logger.Printf("Something went wrong: %s", err)
		return
	}
	logger.Printf("Published message: %s", body)
}
