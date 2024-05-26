package ussd

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/logs"
)

// MenuFunc returns the menu function
type MenuFunc func(params map[string]string) string

// Step a ussd step
type Step struct {
	Run      MenuFunc
	Menu     string
	params   map[string]string
	End      bool
	key      string
	children []*Step
}

// NewStep creates a new step
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

// AddStep adds a new step
func (s *Step) AddStep(key string, child *Step) {
	child.key = key
	s.children = append(s.children, child)
}

// GetResponse returns the response
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

// getMatchingChild returns the matching child
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

// walk goes through the menu and finds matching children
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

// parse takes a ussd string and returns the right child
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
