package ussd

import (
	"fmt"
	"strings"

	"github.com/ochom/gutils/logs"
)

var mainMenu *Step

func New(step *Step) {
	if mainMenu == nil {
		mainMenu = step
	}
}

type Params struct {
	Text        string
	SessionId   string
	PhoneNumber string
}

// Process processes the ussd string and returns the next step
func Process(data Params) (*Step, error) {
	if mainMenu == nil {
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

	mainMenu.params = params
	step := mainMenu.parse(mainMenu.params, parts)
	if step == nil {
		logs.Error("Could not get the correct child for the ussd string")
		return nil, fmt.Errorf("could not get the correct child for the ussd string")
	}

	return step, nil
}
