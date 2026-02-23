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

// StreamSdk provides real-time event streaming capabilities via an external service.
//
// Configuration via environment variables:
//   - STREAMX_URL: API endpoint URL (default: https://api.StreamSdk.co.ke)
//   - STREAMX_API_KEY: Authentication API key
//   - STREAMX_INSTANCE_ID: Instance identifier (default: "default")
type StreamSdk struct {
	url        string
	apiKey     string
	instanceID string
}

// PublishStream sends an event to a channel for real-time distribution.
//
// Parameters:
//   - channel: The channel name subscribers listen to
//   - event: The event type/name
//   - data: The event payload (will be JSON encoded)
//
// Example:
//
//	sdk := pubsub.NewStreamX()
//
//	// Publish a notification
//	sdk.PublishStream("user-123", "notification", map[string]string{
//		"title": "New Message",
//		"body":  "You have a new message",
//	})
//
//	// Publish an order update
//	sdk.PublishStream("orders", "order.updated", order)
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

// NewStreamX creates a new StreamSdk instance for real-time event publishing.
//
// Optional parameters (in order):
//   - params[0]: Instance ID (overrides STREAMX_INSTANCE_ID)
//   - params[1]: API URL (overrides STREAMX_URL)
//   - params[2]: API Key (overrides STREAMX_API_KEY)
//
// Example:
//
//	// Using environment variables
//	sdk := pubsub.NewStreamX()
//
//	// With custom instance ID
//	sdk := pubsub.NewStreamX("my-app-instance")
//
//	// With all custom parameters
//	sdk := pubsub.NewStreamX("instance-1", "https://stream.myapp.com", "api-key-123")
//
//	// Publish events
//	sdk.PublishStream("notifications", "new-message", data)
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
