package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// publisher ...
type publisher struct {
	connectionName string
	url            string
	exchange       string
	queue          string
	exchangeType   ExchangeType
	routingKey     string
}

// NewPublisher  creates a new publisher to rabbit
func NewPublisher(rabbitURL, exchangeName string) Publisher {
	return &publisher{
		url:          rabbitURL,
		exchange:     exchangeName,
		exchangeType: Direct,
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

	defer conn.Close()

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
