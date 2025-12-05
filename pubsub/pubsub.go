package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// ExchangeType ...
type ExchangeType string

var (
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
	FanOut  ExchangeType = "fanout"
	Headers ExchangeType = "headers"
	Delayed ExchangeType = "x-delayed-message"
)

type Publisher interface {
	SetConnectionName(string)
	SetExchangeType(ExchangeType)
	SetQueueName(string)
	SetRoutingKey(string)
	Publish([]byte) error
	PublishWithDelay([]byte, time.Duration) error
}

type Consumer interface {
	SetConnectionName(string)
	SetExchangeName(string)
	SetRoutingKey(string)
	SetDurable(bool)
	SetDeleteWhenUnused(bool)
	SetTag(string)
	SetAutoAck(bool)
	SetExclusive(bool)
	SetNoLocal(bool)
	SetNoWait(bool)
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
