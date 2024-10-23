package pubsub

import (
	"fmt"

	"github.com/streadway/amqp"
)

// consumer ...
type consumer struct {
	connectionName string
	url            string
	queue          string

	// more
	tag       string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
}

// SetConnectionName implements Consumer.
func (c *consumer) SetConnectionName(connectionName string) {
	c.connectionName = connectionName
}

// SetAutoAck implements Consumer.
func (c *consumer) SetAutoAck(autoAck bool) {
	c.autoAck = autoAck
}

// SetExclusive implements Consumer.
func (c *consumer) SetExclusive(exclusive bool) {
	c.exclusive = exclusive
}

// SetNoLocal implements Consumer.
func (c *consumer) SetNoLocal(noLocal bool) {
	c.noLocal = noLocal
}

// SetNoWait implements Consumer.
func (c *consumer) SetNoWait(noWait bool) {
	c.noWait = noWait
}

// SetTag implements Consumer.
func (c *consumer) SetTag(tag string) {
	c.tag = tag
}

// Create a new consumer instance
func NewConsumer(rabbitURL, queueName string) Consumer {
	return &consumer{
		url:       rabbitURL,
		queue:     queueName,
		autoAck:   true,
		exclusive: false,
		noLocal:   false,
		noWait:    false,
	}
}

// Consume consume messages from the channels
func (c *consumer) Consume(workerFunc func([]byte)) error {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": c.connectionName,
		},
	}

	conn, err := amqp.DialConfig(c.url, cfg)
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
		q.Name,      // queue
		c.tag,       // consumerTag
		c.autoAck,   // auto-ack
		c.exclusive, // exclusive
		c.noLocal,   // no-local
		c.noWait,    // no-wait
		nil,         // args
	)

	if err != nil {
		return fmt.Errorf("failed to consume messages: %s", err.Error())
	}

	for message := range deliveries {
		workerFunc(message.Body)
	}

	return nil
}
