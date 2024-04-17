package pubsub

import "fmt"

// Consumer ...
type Consumer struct {
	url      string
	exchange string
	queue    string
}

// Create a new consumer instance
func NewConsumer(rabbitURL, exchange, queue string) *Consumer {
	return &Consumer{rabbitURL, exchange, queue}
}

// Consume consume messages from the channels
func (c *Consumer) Consume(workerFunc func([]byte)) error {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}
	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	deliveries, err := ch.Consume(
		c.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return fmt.Errorf("failed to consume messages: %s", err.Error())
	}

	for d := range deliveries {
		workerFunc(d.Body)
	}

	return nil
}
