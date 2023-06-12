package pubsub

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// publisher ...
type publisher struct {
	url      string
	exchange string
	queue    string
}

// newPublisher ...
func newPublisher(queueName string) *publisher {
	exchange := fmt.Sprintf("%s-exchange", queueName)
	return &publisher{rabbitURL, exchange, queueName}
}

// publish ...
func (p *publisher) publish(body []byte, delay time.Duration, delayed bool) error {
	conn, ch, err := initQ(p.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	if err := bind(ch, p.exchange, p.queue, delayed); err != nil {
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
func PublishWithDelay(queueName string, body []byte, delay time.Duration) error {
	p := newPublisher(fmt.Sprintf("%s-%s", queuePrefix, queueName))
	return p.publish(body, delay, true)
}

// Publish ...
func Publish(queueName string, body []byte) error {
	p := newPublisher(fmt.Sprintf("%s-%s", queuePrefix, queueName))
	return p.publish(body, 0, false)
}
