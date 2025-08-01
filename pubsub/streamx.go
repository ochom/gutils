package pubsub

import (
	"fmt"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/jsonx"
	"github.com/ochom/gutils/logs"
)

var (
	url        = env.Get("STREAMX_URL", "https://api.StreamSdk.co.ke")
	apiKey     = env.Get("STREAMX_API_KEY", "")
	instanceID = env.Get("STREAMX_INSTANCE_ID", "default")
)

type StreamSdk struct {
	url        string
	apiKey     string
	instanceID string
}

// PublishStream publishes a message to the stream
func (s StreamSdk) PublishStream(channel string, event string, data any) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": s.apiKey,
	}

	url := fmt.Sprintf("%s/publish", s.url)
	res, err := gttp.Post(url, headers, jsonx.Encode(map[string]any{
		"instance_id": s.instanceID,
		"channel":     channel,
		"event":       event,
		"data":        data,
	}))

	if err != nil {
		logs.Error("Failed to publish message to stream: %v", err)
		return
	}

	if res.StatusCode != 200 {
		logs.Error("Failed to publish message to stream: %v", string(res.Body))
		return
	}
}

// NewStreamX create new instance of StreamSdk
// with optional parameters for instance ID, URL, and API key.
func NewStreamX(params ...string) (sdk *StreamSdk) {
	sdk = &StreamSdk{
		url:        url,
		apiKey:     apiKey,
		instanceID: instanceID,
	}

	if len(params) > 0 {
		sdk.instanceID = params[0]
	}

	if len(params) > 1 {
		sdk.url = params[1]
	}

	if len(params) > 2 {
		sdk.apiKey = params[2]
	}

	return sdk
}
