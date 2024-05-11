package ussd

import (
	"fmt"
	"strings"

	"github.com/ochom/gutils/logs"
)

var mainMenu *Step

func InitMenu(step *Step) {
	if mainMenu == nil {
		mainMenu = step
	}
}

type Params struct {
	Ussd        string
	SessionId   string
	PhoneNumber string
}

func Parse(data Params) (*Step, error) {
	if mainMenu == nil {
		return nil, fmt.Errorf("mainMenu has not been created")
	}

	data.Ussd = strings.TrimSuffix(data.Ussd, "#")
	data.Ussd = strings.TrimPrefix(data.Ussd, "*")

	parts := []string{}
	if data.Ussd != "" {
		parts = strings.Split(data.Ussd, "*")
	}

	// check if session exists
	if !HasSession(data.SessionId) {
		SetSession(data.SessionId, make(map[string]string))
	}

	mainMenu.Params = map[string]string{
		"phone_number": data.PhoneNumber,
		"session_id":   data.SessionId,
		"ussd":         data.Ussd,
	}
	step := mainMenu.Parse(mainMenu.Params, parts)
	if step == nil {
		logs.Error("Could not get the correct child for the ussd string")
		return nil, fmt.Errorf("could not get the correct child for the ussd string")
	}

	return step, nil
}
