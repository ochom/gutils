package ussd

import (
	"fmt"
	"strings"
)

var root *Step

func New(step *Step) {
	root = step
}

// Parse processes the ussd string and returns the next step
func Parse(data Params) (*Step, error) {
	if root == nil {
		return nil, fmt.Errorf("mainMenu has not been created")
	}

	data.Text = strings.TrimSuffix(data.Text, "#")
	data.Text = strings.TrimPrefix(data.Text, "*")

	parts := []string{}
	if data.Text != "" {
		parts = strings.Split(data.Text, "*")
	}

	params := map[string]string{
		"phone_number": data.PhoneNumber,
		"session_id":   data.SessionId,
		"text":         data.Text,
	}

	if len(parts) > 0 {
		params["input"] = parts[len(parts)-1]
	}

	for k, v := range params {
		SetSession(data.SessionId, k, v)
	}

	step := root.parse(params, parts)
	if step == nil {
		return nil, fmt.Errorf("step not found")
	}

	return step, nil
}
