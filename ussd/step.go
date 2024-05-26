package ussd

import (
	"regexp"
	"strings"
)

// MenuFunc returns the menu function
type MenuFunc func(params map[string]string) string

// Children returns the children
type Children map[string]*Step

// Step a ussd step
type Step struct {
	Menu   MenuFunc
	End    bool
	params map[string]string
	Children
}

// NewStep creates a new step
func NewStep(menuFunc MenuFunc) Step {
	return Step{
		Menu:   menuFunc,
		params: make(map[string]string),
	}
}

// AddStep adds a new step
func (s *Step) AddStep(key string, newStep *Step) {
	if len(s.Children) == 0 {
		s.Children = make(map[string]*Step)
	}

	s.Children[key] = newStep
}

// GetResponse returns the response
func (s *Step) GetResponse() string {
	response := s.Menu(s.params)
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
	child, ok := s.Children[key]
	if ok {
		return child
	}

	// get regex match
	for key := range s.Children {
		if key == "" {
			key = "*"
		}

		matcher, err := regexp.Compile(key)
		if err != nil {
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
func (s *Step) parse(sessionData map[string]string, ussdParts []string) *Step {
	child := s.walk(ussdParts)
	if child == nil {
		return s
	}

	remainingPieces := len(ussdParts)
	if remainingPieces > 1 {
		sessionData["input"] = ussdParts[remainingPieces-1]
	}

	child.params = sessionData
	return child
}
