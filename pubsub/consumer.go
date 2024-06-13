package pubsub

import (
	"fmt"
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
func (c *Consumer) Consume(stop chan bool, workerFunc func([]byte)) error {
	ps, err := initQ(c.url)
	if err != nil {
		return fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}
	defer ps.Close()

	if err := initPubSub(ps.ch, c.exchange, c.queue, string(c.config.Type)); err != nil {
		return fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	deliveries, err := ps.ch.Consume(
		c.queue,            // queue
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

	for {
		select {
		case <-stop:
			return fmt.Errorf("stop signal received")
		case message := <-deliveries:
			workerFunc(message.Body)
		}
	}
}
