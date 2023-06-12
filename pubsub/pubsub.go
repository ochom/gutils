package pubsub

import (
	"fmt"

	"github.com/ochom/gutils/helpers"
	"github.com/streadway/amqp"
)

var (
	// rabbitURL is the URL of the RabbitMQ server
	rabbitURL = helpers.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")

	// queuePrefix used to prefix the queue name to avoid conflict with other services
	queuePrefix = helpers.GetEnv("QUEUE_PREFIX", "dev")
)

func initQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = ch.Qos(1, 0, false) // fair dispatch
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// initPubSub ...
func initPubSub(ch *amqp.Channel, args amqp.Table, channelType, exchangeName, queueName string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		channelType,  // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		args,         // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
	}

	// declare queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		q.Name,       // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}

// bind ...
func bind(ch *amqp.Channel, exchangeName, queueName string, delayed bool) error {
	args := make(amqp.Table)
	if delayed {
		args["x-delayed-type"] = "direct"
		channelType := "x-delayed-message"
		return initPubSub(ch, args, channelType, exchangeName, queueName)
	}

	args["x-queue-mode"] = "lazy"
	channelType := "direct"
	return initPubSub(ch, args, channelType, exchangeName, queueName)
}
