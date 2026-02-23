package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// publisher implements the Publisher interface for RabbitMQ message publishing.
type publisher struct {
	connectionName string
	url            string
	exchange       string
	queue          string
	exchangeType   ExchangeType
	routingKey     string
}

// NewPublisher creates a new publisher for sending messages to RabbitMQ.
//
// Parameters:
//   - rabbitURL: AMQP connection URL (e.g., "amqp://guest:guest@localhost:5672/")
//   - exchangeName: Name of the exchange to publish to
//   - queueName: Name of the queue (used for declaring and binding)
//
// Example:
//
//	// Create a publisher
//	publisher := pubsub.NewPublisher(
//		"amqp://guest:guest@localhost:5672/",
//		"orders",      // exchange
//		"order-queue", // queue
//	)
//
//	// Configure and publish
//	publisher.SetRoutingKey("order.created")
//	err := publisher.Publish(jsonx.Encode(order))
//
//	// Publish with delay
//	publisher.SetExchangeType(pubsub.Delayed)
//	err := publisher.PublishWithDelay(jsonx.Encode(reminder), 24*time.Hour)
func NewPublisher(rabbitURL, exchangeName, queueName string) Publisher {
	return &publisher{
		url:          rabbitURL,
		exchange:     exchangeName,
		exchangeType: Direct,
		queue:        queueName,
	}
}

// SetQueueName ...
func (p *publisher) SetQueueName(queueName string) {
	p.queue = queueName
}

// SetConnectionName ...
func (p *publisher) SetConnectionName(connectionName string) {
	p.connectionName = connectionName
}

// SetRoutingKey ...
func (p *publisher) SetRoutingKey(routingKey string) {
	p.routingKey = routingKey
}

// SetExchangeType ...
func (p *publisher) SetExchangeType(exchangeType ExchangeType) {
	p.exchangeType = exchangeType
}

// PublishWithDelay ...
func (p *publisher) PublishWithDelay(body []byte, delay time.Duration) error {
	return p.publish(body, delay)
}

// Publish ...
func (p *publisher) Publish(body []byte) error {
	return p.publish(body, 0)
}

// declare create exchange and queue
func (p *publisher) declare(ch *amqp.Channel) error {
	return declare(ch, p.exchange, p.queue, p.routingKey, p.exchangeType)
}

// publish ...
func (p *publisher) publish(body []byte, delay time.Duration) error {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": p.connectionName,
		},
	}

	conn, err := amqp.DialConfig(p.url, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	defer func() {
		_ = conn.Close()
	}()

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	err = channel.Qos(10, 0, false) // fair dispatch
	if err != nil {
		return fmt.Errorf("failed to set QoS: %s", err.Error())
	}

	if err := p.declare(channel); err != nil {
		return fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	// publish message to exchange
	err = channel.Publish(
		p.exchange,   // exchange
		p.routingKey, // routing key
		true,         // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers: map[string]any{
				"x-delay": delay.Milliseconds(),
			},
		},
	)

	return err
}
