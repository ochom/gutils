package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// publisher ...
type publisher struct {
	url          string
	exchange     string
	queue        string
	exchangeType string
}

// NewPublisher  creates a new publisher to rabbit
func NewPublisher(rabbitURL, exchange, queue string, exchangeType ExchangeType) *publisher {
	return &publisher{rabbitURL, exchange, queue, string(exchangeType)}
}

// publish ...
func (p *publisher) publish(body []byte, delay time.Duration) error {
	ps, err := initQ(p.url)
	if err != nil {
		return fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}
	defer ps.Close()

	if err := initPubSub(ps.ch, p.exchange, p.queue, p.exchangeType); err != nil {
		return fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	// publish message to exchange
	err = ps.ch.Publish(
		p.exchange, // exchange
		p.queue,    // routing key
		true,       // mandatory
		false,      // immediate
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

// PublishWithDelay ...
func (p *publisher) PublishWithDelay(body []byte, delay time.Duration) error {
	return p.publish(body, delay)
}

// Publish ...
func (p *publisher) Publish(body []byte) error {
	return p.publish(body, 0)
}
