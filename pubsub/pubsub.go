package pubsub

import "time"

// ExchangeType ...
type ExchangeType string

var (
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
	FanOut  ExchangeType = "fanout"
	Headers ExchangeType = "headers"
	Delayed ExchangeType = "x-delayed-message"
)

type Publisher interface {
	SetConnectionName(string)
	SetExchangeType(ExchangeType)
	SetRoutingKey(string)
	Publish([]byte) error
	PublishWithDelay([]byte, time.Duration) error
}

type Consumer interface {
	SetConnectionName(string)
	SetTag(string)
	SetAutoAck(bool)
	SetExclusive(bool)
	SetNoLocal(bool)
	SetNoWait(bool)
	Consume(func([]byte)) error
}
