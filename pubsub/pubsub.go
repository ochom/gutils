// Package pubsub provides RabbitMQ publish/subscribe functionality with support
// for various exchange types including delayed messaging.
//
// This package simplifies working with RabbitMQ by providing a builder-pattern
// interface for configuring publishers and consumers.
//
// Supported exchange types:
//   - Direct: Routes messages to queues based on exact routing key match
//   - Topic: Routes messages based on pattern matching in routing keys
//   - FanOut: Broadcasts messages to all bound queues
//   - Headers: Routes based on message headers
//   - Delayed: Supports delayed message delivery using x-delayed-message plugin
//
// Example publisher:
//
//	publisher := pubsub.NewPublisher(
//		\"amqp://guest:guest@localhost:5672/\",
//		\"notifications\",  // exchange name
//		\"email-queue\",    // queue name
//	)
//	publisher.SetRoutingKey(\"email.send\")
//
//	err := publisher.Publish(jsonx.Encode(emailData))
//	if err != nil {
//		log.Error(\"Failed to publish: %v\", err)
//	}
//
// Example consumer:
//
//	consumer := pubsub.NewConsumer(
//		\"amqp://guest:guest@localhost:5672/\",
//		\"email-queue\",
//	)
//
//	err := consumer.Consume(func(msg amqp.Delivery) {
//		var email EmailData
//		json.Unmarshal(msg.Body, &email)
//		sendEmail(email)
//	})
package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// ExchangeType represents the type of RabbitMQ exchange.
type ExchangeType string

// Supported exchange types
var (
	// Direct exchange routes messages to queues based on exact routing key match
	Direct ExchangeType = "direct"
	// Topic exchange routes messages based on wildcard pattern matching in routing keys
	Topic ExchangeType = "topic"
	// FanOut exchange broadcasts messages to all bound queues, ignoring routing keys
	FanOut ExchangeType = "fanout"
	// Headers exchange routes messages based on message header attributes
	Headers ExchangeType = "headers"
	// Delayed exchange supports delayed message delivery (requires RabbitMQ plugin)
	Delayed ExchangeType = "x-delayed-message"
)

// Publisher defines the interface for publishing messages to RabbitMQ.
type Publisher interface {
	// SetConnectionName sets a name for the connection (visible in RabbitMQ management)
	SetConnectionName(string)
	// SetExchangeType sets the exchange type (direct, topic, fanout, headers, delayed)
	SetExchangeType(ExchangeType)
	// SetQueueName sets the target queue name
	SetQueueName(string)
	// SetRoutingKey sets the routing key for message routing
	SetRoutingKey(string)
	// Publish sends a message immediately
	Publish([]byte) error
	// PublishWithDelay sends a message with a delay (requires delayed exchange)
	PublishWithDelay([]byte, time.Duration) error
}

// Consumer defines the interface for consuming messages from RabbitMQ.
type Consumer interface {
	// SetConnectionName sets a name for the connection
	SetConnectionName(string)
	// SetExchangeName sets the exchange to bind to
	SetExchangeName(string)
	// SetRoutingKey sets the routing key for queue binding
	SetRoutingKey(string)
	// SetDurable sets whether the queue survives broker restart
	SetDurable(bool)
	// SetDeleteWhenUnused sets whether the queue is deleted when no consumers
	SetDeleteWhenUnused(bool)
	// SetTag sets the consumer tag for identification
	SetTag(string)
	// SetAutoAck sets whether messages are auto-acknowledged
	SetAutoAck(bool)
	// SetExclusive sets whether this is an exclusive consumer
	SetExclusive(bool)
	// SetNoLocal sets whether to not receive messages published on this connection
	SetNoLocal(bool)
	// SetNoWait sets whether to wait for server confirmation
	SetNoWait(bool)
	// Consume starts consuming messages, calling workerFunc for each message
	Consume(func(amqp.Delivery)) error
}

// declare create exchange and queue
func declare(ch *amqp.Channel, exchange, queue, routingKey string, exchangeType ExchangeType) error {
	err := ch.ExchangeDeclare(
		exchange,             // name
		string(exchangeType), // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
	}

	// if queue is empty, then no need to declare queue
	if queue == "" {
		return nil
	}

	// declare queue
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	return bindQueue(ch, exchange, q.Name, routingKey)
}

// bind queue to exchange ...
func bindQueue(ch *amqp.Channel, exchange, queue, routingKey string) error {
	if exchange == "" || queue == "" {
		return nil
	}

	err := ch.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}
