package pubsub

import (
	"time"

	"github.com/streadway/amqp"
)

// publisher ...
type publisher struct {
	url      string
	exchange string
	queue    string
}

// NewPublisher  creates a new publisher to rabbit
func NewPublisher(rabbitURL, exchange, queue string) *publisher {
	return &publisher{rabbitURL, exchange, queue}
}

// publish ...
func (p *publisher) publish(body []byte, delay time.Duration) error {
	conn, ch, err := initQ(p.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, p.exchange, p.queue); err != nil {
		return err
	}

	// publish message to exchange
	err = ch.Publish(
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
