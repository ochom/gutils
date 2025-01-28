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
func NewPublisher(rabbitURL, exchangeName, queueName string) Publisher {
	return &publisher{
		url:          rabbitURL,
		exchange:     exchangeName,
		queue:        queueName,
		exchangeType: Direct,
	}
}

func (p *publisher) SetConnectionName(connectionName string) {
	p.connectionName = connectionName
}

func (p *publisher) SetRoutingKey(routingKey string) {
	p.routingKey = routingKey
}

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

// initPubSub ...
func (p *publisher) initPubSub(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		p.exchange,             // name
		string(p.exchangeType), // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
	}

	// declare queue
	q, err := ch.QueueDeclare(
		p.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		p.routingKey, // routing key
		p.exchange,   // exchange
		false,        // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
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

	if err := p.initPubSub(channel); err != nil {
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
