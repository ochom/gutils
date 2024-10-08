package pubsub

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Consumer ...
type Consumer struct {
	url      string
	exchange string
	queue    string
	config   Config
}

type Config struct {
	Type      ExchangeType
	Tag       string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

var defaultConfig = Config{
	Type:      Direct,
	Tag:       "",
	AutoAck:   true,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
}

// Create a new consumer instance
func NewConsumer(rabbitURL, exchange, queue string, config ...Config) *Consumer {
	c := Consumer{url: rabbitURL, exchange: exchange, queue: queue, config: defaultConfig}

	if len(config) > 0 {
		c.config = config[0]
	}

	return &c
}

// Consume consume messages from the channels
func (c *Consumer) Consume(workerFunc func([]byte)) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	deliveries, err := ch.Consume(
		q.Name,             // queue
		c.config.Tag,       // consumerTag
		c.config.AutoAck,   // auto-ack
		c.config.Exclusive, // exclusive
		c.config.NoLocal,   // no-local
		c.config.NoWait,    // no-wait
		nil,                // args
	)

	if err != nil {
		return fmt.Errorf("failed to consume messages: %s", err.Error())
	}

	for message := range deliveries {
		workerFunc(message.Body)
	}

	return nil
}
