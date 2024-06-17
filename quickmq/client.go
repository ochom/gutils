package quickmq

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/r3labs/sse/v2"
)

// Client ...
type Client struct {
	url      string
	username string
	password string
	queue    string
}

// NewClient  creates a new  quickmq client
func NewClient(quickUrl, queue string) *Client {
	url, username, password, err := parseUrl(quickUrl)
	if err != nil {
		panic(err)
	}

	return &Client{url: url, username: username, password: password, queue: queue}
}

// publish ...
func (p *Client) publish(body []byte, delay time.Duration) error {
	message := map[string]any{
		"body":  string(body),
		"delay": delay,
		"queue": p.queue,
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.username, p.password))),
	}

	url := fmt.Sprintf("%s/publish", p.url)
	res, err := gttp.Post(url, headers, helpers.ToBytes(message), time.Minute)
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err.Error())
	}

	if res.Status != 200 {
		return fmt.Errorf("failed to publish message: %s", string(res.Body))
	}

	return nil
}

// PublishWithDelay ...
func (p *Client) PublishWithDelay(body []byte, delay time.Duration) error {
	return p.publish(body, delay)
}

// Publish ...
func (p *Client) Publish(body []byte) error {
	return p.publish(body, 0)
}

// Consume consume messages from the channels
func (c *Client) Consume(stop chan bool, workerFunc func([]byte)) error {
	events := make(chan *sse.Event)
	url := fmt.Sprintf("%s/subscribe?queue=%s", c.url, c.queue)

	client := sse.NewClient(url, func(sseClient *sse.Client) {
		headers := map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password))),
		}
		sseClient.Headers = headers
	})
	if err := client.SubscribeChanRaw(events); err != nil {
		return err
	}

	for {
		select {
		case <-stop:
			return fmt.Errorf("stop signal received")
		case message := <-events:
			if bytes.Equal(message.Data, []byte(`{}`)) {
				continue
			}
			workerFunc(message.Data)
		}
	}
}
