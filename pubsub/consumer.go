package pubsub

import (
	"fmt"

	"github.com/streadway/amqp"
)

// consumer ...
type consumer struct {
	connectionName string
	url            string
	exchange       string
	queue          string
	routingKey     string

	// basic
	durable          bool
	deleteWhenUnused bool

	// more
	tag       string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
}

// SetExchangeName implements Consumer.
func (c *consumer) SetExchangeName(exchangeName string) {
	c.exchange = exchangeName
}

// SetRoutingKey implements Consumer.
func (c *consumer) SetRoutingKey(routingKey string) {
	c.routingKey = routingKey
}

// SetConnectionName implements Consumer.
func (c *consumer) SetConnectionName(connectionName string) {
	c.connectionName = connectionName
}

// SetDurable implements Consumer.
func (c *consumer) SetDurable(durable bool) {
	c.durable = durable
}

// SetDeleteWhenUnused implements Consumer.
func (c *consumer) SetDeleteWhenUnused(deleteWhenUnused bool) {
	c.deleteWhenUnused = deleteWhenUnused
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
		url:     rabbitURL,
		queue:   queueName,
		durable: true,
		autoAck: true,
	}
}

// Consume consume messages from the channels
func (c *consumer) Consume(workerFunc func(amqp.Delivery)) error {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": c.connectionName,
		},
	}

	conn, err := amqp.DialConfig(c.url, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	defer func() {
		_ = ch.Close()
	}()

	q, err := ch.QueueDeclare(
		c.queue,            // name
		c.durable,          // durable
		c.deleteWhenUnused, // delete when unused
		c.exclusive,        // exclusive
		c.noWait,           // no-wait
		nil,                // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	err = bindQueue(ch, c.exchange, q.Name, c.routingKey)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
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
		workerFunc(message)
	}

	return nil
}
