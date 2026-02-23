package ussd

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/logs"
)

// MenuFunc is a function type that generates dynamic menu content.
// It receives session parameters and returns the menu text to display.
//
// Example:
//
//	balanceMenu := ussd.NewStep(func(params map[string]string) string {
//		phone := params["phone_number"]
//		balance := getBalance(phone)
//		return fmt.Sprintf("Your balance is: KES %s\\n0. Back to Menu", balance)
//	})
type MenuFunc func(params map[string]string) string

// Step represents a single step/screen in the USSD menu flow.
// Each step can have child steps for nested menus.
type Step struct {
	// Run is a function that generates dynamic content (mutually exclusive with Menu)
	Run MenuFunc
	// Menu is static menu text (mutually exclusive with Run)
	Menu string
	// params contains session data passed to the step
	params map[string]string
	// End marks this as a terminal step (session will be cleared after response)
	End bool
	// key is the input that leads to this step
	key string
	// children are the sub-steps accessible from this step
	children []*Step
}

// NewStep creates a new USSD menu step.
// Accepts either a static string or a MenuFunc for dynamic content.
//
// Example:
//
//	// Static menu
//	mainMenu := ussd.NewStep("Welcome!\\n1. Balance\\n2. Transfer")
//
//	// Dynamic menu with function
//	balanceMenu := ussd.NewStep(func(params map[string]string) string {
//		balance := checkBalance(params["phone_number"])
//		return fmt.Sprintf("Your balance: KES %s", balance)
//	})
func NewStep[T string | MenuFunc](menu T) *Step {
	step := &Step{
		params: make(map[string]string),
	}

	ok := reflect.TypeOf(menu).Kind() == reflect.Func
	if ok {
		step.Run = any(menu).(MenuFunc)
	} else {
		step.Menu = any(menu).(string)
	}

	return step
}

// AddStep adds a child step accessible by the given key.
// The key is the user input that will navigate to this step.
//
// Keys can be:
//   - Exact match: "1", "2", "0"
//   - Wildcard: "*" matches any input
//   - Regex pattern: "[0-9]+" matches numeric input
//
// Example:
//
//	mainMenu := ussd.NewStep("Main Menu\n1. Balance\n2. Transfer")
//
//	// Exact match
//	mainMenu.AddStep("1", balanceStep)
//	mainMenu.AddStep("2", transferStep)
//
//	// Wildcard for "back" handling
//	mainMenu.AddStep("0", backStep)
//
//	// Any input (for text input)
//	inputStep := ussd.NewStep("Enter amount:")
//	inputStep.AddStep("*", confirmationStep) // Matches any input
func (s *Step) AddStep(key string, child *Step) {
	child.key = key
	s.children = append(s.children, child)
}

// GetResponse generates the USSD response string for this step.
// Automatically prefixes with "CON " (continue) or "END " (terminate) based on the step's End flag.
//
// Response prefixes:
//   - "CON ": User can continue interacting (session stays active)
//   - "END ": Final message (session will be terminated)
//
// If the menu text already starts with "CON " or "END ", it's returned as-is.
//
// Example:
//
//	// Continuing step
//	step := ussd.NewStep("Enter PIN:")
//	step.GetResponse() // Returns: "CON Enter PIN:"
//
//	// Terminal step
//	step := ussd.NewStep("Thank you!")
//	step.End = true
//	step.GetResponse() // Returns: "END Thank you!"
//
//	// Pre-formatted response
//	step := ussd.NewStep(func(p map[string]string) string {
//		return "END Your transaction is complete"
//	})
//	step.GetResponse() // Returns: "END Your transaction is complete"
func (s *Step) GetResponse() string {
	response := s.Menu
	if s.Menu == "" {
		response = s.Run(s.params)
	}

	if strings.HasPrefix(response, "END ") {
		RemoveSession(s.params["session_id"])
		return response
	}

	if strings.HasPrefix(response, "CON ") {
		return response
	}

	if s.End {
		RemoveSession(s.params["session_id"])
		return "END " + response
	}

	return "CON " + response
}

// getMatchingChild finds a child step that matches the given input key.
// First checks for exact matches, then falls back to regex pattern matching.
func (s *Step) getMatchingChild(key string) *Step {
	// get exact match
	child := arrays.Find(s.children, func(child *Step) bool {
		return strings.EqualFold(child.key, key)
	})

	if child != nil {
		return child
	}

	for _, child := range s.children {
		if child.key == "" || child.key == "*" {
			child.key = ".*"
		}

		matcher, err := regexp.Compile(child.key)
		if err != nil {
			logs.Error("could not compile regexp: %v", err)
			continue
		}

		if matcher.MatchString(key) {
			return child
		}
	}

	return nil
}

// walk traverses the menu tree based on the USSD input parts.
// Returns the step that corresponds to the complete navigation path.
func (s *Step) walk(ussdParts []string) *Step {
	remainingPieces := len(ussdParts)

	// break if no piece is left
	if remainingPieces == 0 {
		return s
	}

	// break if only once piece is left by finding the child that matches the remaining piece
	if remainingPieces == 1 {
		return s.getMatchingChild(ussdParts[0])
	}

	// get the child that matches the first piece
	if child := s.getMatchingChild(ussdParts[0]); child != nil {
		return child.walk(ussdParts[1:])
	}

	return nil
}

// parse processes the USSD parts and returns the appropriate step with parameters.
func (s *Step) parse(params map[string]string, ussdParts []string) *Step {
	step := s.walk(ussdParts)
	if step == nil {
		logs.Warn("parts %v produced no step", ussdParts)
		return nil
	}

	remainingPieces := len(ussdParts)
	if remainingPieces > 1 {
		params["input"] = ussdParts[remainingPieces-1]
	}

	step.params = params
	return step
}
