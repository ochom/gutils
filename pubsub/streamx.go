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
	Data       any    `json:"data"`
}

type StreamSdk struct {
	url    string
	apiKey string
}

// PublishStream publishes a message to the stream
func (s StreamSdk) PublishStream(message *StreamMessage) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": s.apiKey,
	}

	url := fmt.Sprintf("%s/publish", s.url)
	res, err := gttp.Post(url, headers, helpers.ToBytes(message))
	if err != nil {
		logs.Error("Failed to publish message to stream: %v", err)
		return
	}

	if res.StatusCode != 200 {
		logs.Error("Failed to publish message to stream: %v", string(res.Body))
		return
	}
}

type StreamSdkConfig struct {
	Url    string
	ApiKey string
}

var DefaultConfig = &StreamSdkConfig{
	Url:    env.Get("STREAMX_URL", "https://api.StreamSdk.co.ke"),
	ApiKey: env.Get("STREAMX_API_KEY"),
}

func NewStreamX(cfg *StreamSdkConfig) (sdk *StreamSdk) {
	sdk = &StreamSdk{
		url:    DefaultConfig.Url,
		apiKey: DefaultConfig.ApiKey,
	}

	if cfg == nil {
		return
	}

	if cfg.Url != "" {
		sdk.url = cfg.Url
	}

	if cfg.ApiKey != "" {
		sdk.apiKey = cfg.ApiKey
	}

	return sdk
}
