package pubsub

import (
	"fmt"

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

// Pubsub ...
type Pubsub struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Close ...
func (p *Pubsub) Close() {
	if err := p.conn.Close(); err != nil {
		fmt.Println("failed to close connection: ", err.Error())
	}

	if err := p.ch.Close(); err != nil {
		fmt.Println("failed to close channel: ", err.Error())
	}
}

func initQ(url string) (*Pubsub, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	err = ch.Qos(10, 0, false) // fair dispatch
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %s", err.Error())
	}

	return &Pubsub{conn: conn, ch: ch}, nil
}

// initPubSub ...
func initPubSub(ch *amqp.Channel, exchangeName, queueName, exchangeType string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
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
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		q.Name,       // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}
