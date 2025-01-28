package pubsub

import (
	"fmt"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// StreamMessage ...
type StreamMessage struct {
	InstanceID string `json:"instanceID"`
	Channel    string `json:"channel"`
	ID         string `json:"id"`
	Event      string `json:"event"`
	Data       any    `json:"message"`
}

type StreamX struct {
	Url    string
	APIKey string
}

var streamX *StreamX

func InitStreamX(apiKey string) {
	streamX = &StreamX{APIKey: apiKey}
}

func (s *StreamX) publish(message *StreamMessage) {
	if streamX == nil {
		logs.Error("StreamX not initialized")
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": env.Get("STREAMX_API_KEY"),
	}

	url := fmt.Sprintf("%s/publish", env.Get("STREAMX_URL", "https://api.streamx.co.ke"))
	res, err := gttp.Post(url, headers, helpers.ToBytes(message))
	if err != nil {
		logs.Error("Failed to publish message to stream: %v", err)
		return
	}

	if res.StatusCode != 200 {
		logs.Error("Failed to publish message to stream: %v", res.Body)
	}

	logs.Info("StreamMessage published to StreamX ==> msgID: %s", message.ID)
}

// PublishStream publishes a message to the stream
func PublishStream(message *StreamMessage) {
	go streamX.publish(message)
}
