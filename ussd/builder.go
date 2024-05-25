package ussd

import (
	"regexp"
	"strings"

	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/logs"
)

// MenuFunc returns the menu function
type MenuFunc func(params map[string]string) string

type Children []*Step

// Step a ussd step
type Step struct {
	Key      string
	Menu     MenuFunc
	End      bool
	Children Children
	params   map[string]string
}

// NewMenu creates a new step
func NewMenu(menuFunc MenuFunc) Step {
	return Step{
		Menu:   menuFunc,
		params: make(map[string]string),
	}
}

// AddStep adds a new step
func (s *Step) AddStep(step Step) {
	s.Children = append(s.Children, &step)
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
func (s *Step) getMatchingChild(ussdPart string) *Step {
	// first check exact matching children
	child := arrays.Find(s.Children, func(c *Step) bool {
		return strings.EqualFold(ussdPart, c.Key)
	})

	if child != nil {
		return child
	}

	// then check if the key is a wildcard
	child = arrays.Find(s.Children, func(c *Step) bool {
		return c.Key == ""
	})

	return child
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

	// first check kids that exactly match the first piece
	for _, child := range s.Children {
		if strings.EqualFold(ussdParts[0], child.Key) {
			return child.walk(ussdParts[1:])
		}
	}

	// check any item that matches the input as a regex
	for _, child := range s.Children {
		if child.Key == "" {
			child.Key = "*"
		}
		matcher, err := regexp.Compile(child.Key)
		if err != nil {
			continue
		}

		if matcher.MatchString(ussdParts[0]) {
			return child.walk(ussdParts[1:])
		}
	}

	return nil
}

// parse takes a ussd string and returns the right child
func (s *Step) parse(sessionData map[string]string, ussdParts []string) *Step {
	child := s.walk(ussdParts)
	if child == nil {
		logs.Error("Could not get the correct child for the ussd string")
		return nil
	}

	remainingPieces := len(ussdParts)
	if remainingPieces > 1 {
		sessionData["input"] = ussdParts[remainingPieces-1]
	}

	child.params = sessionData
	return child
}
