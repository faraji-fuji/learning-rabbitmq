package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// InitPubSub ...
func InitPubSub(conn *amqp.Connection, exchangeName, queueName string) (ch *amqp.Channel, err error) {

	ch, err = conn.Channel()
	if err != nil {
		return
	}

	err = ch.Qos(1, 0, false) // fair dispatch
	if err != nil {
		return
	}

	err = ch.ExchangeDeclare(
		exchangeName,        // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)
	if err != nil {
		return
	}

	// declare queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		q.Name,       // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return
	}

	return ch, nil
}
