package ussd

import (
	"fmt"
	"strings"
)

var root *Step

// New registers the root menu step for the USSD application.
// Must be called before processing any USSD requests.
//
// Example:
//
//	mainMenu := ussd.NewStep("Welcome\n1. Option A\n2. Option B")
//	mainMenu.AddStep("1", optionAStep)
//	mainMenu.AddStep("2", optionBStep)
//
//	ussd.New(mainMenu)
func New(step *Step) {
	root = step
}

// Parse processes a USSD request and returns the appropriate menu step.
// Extracts session data and navigates the menu tree based on user input.
//
// The Text field in params follows standard USSD format:
//   - Empty string: Show root menu
//   - "1": User selected option 1
//   - "1*2*3": User navigated through options 1 > 2 > 3
//
// Available parameters in the step's params map:
//   - phone_number: The user's phone number
//   - session_id: The current session ID
//   - text: The full input text
//   - input: The last user input
//
// Example:
//
//	func HandleUSSD(w http.ResponseWriter, r *http.Request) {
//		params := ussd.Params{
//			SessionId:   r.FormValue("sessionId"),
//			PhoneNumber: r.FormValue("phoneNumber"),
//			Text:        r.FormValue("text"),
//		}
//
//		step, err := ussd.Parse(params)
//		if err != nil {
//			fmt.Fprint(w, "END An error occurred")
//			return
//		}
//
//		fmt.Fprint(w, step.GetResponse())
//	}
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
