package ussd

// session stores USSD session data keyed by session ID.
// Each session contains key-value pairs for storing user inputs and state.
var session map[string]map[string]string

func init() {
	session = make(map[string]map[string]string)
}

// HasSession checks if a session exists for the given session ID.
//
// Example:
//
//	if ussd.HasSession(sessionId) {
//		// Session exists, restore state
//	} else {
//		// New session, start fresh
//	}
func HasSession(sessionId string) bool {
	_, ok := session[sessionId]
	return ok
}

// GetSession retrieves all session data for a session ID.
// Returns an empty map if the session doesn't exist.
//
// Example:
//
//	data := ussd.GetSession(sessionId)
//	amount := data["amount"]
//	recipient := data["recipient"]
func GetSession(sessionId string) map[string]string {
	val, ok := session[sessionId]
	if !ok {
		return make(map[string]string)
	}

	return val
}

// SetSession stores a key-value pair in the session.
// Creates the session if it doesn't exist.
//
// Example:
//
//	// Store user's menu selection
//	ussd.SetSession(sessionId, "menu_selection", "1")
//
//	// Store transfer amount
//	ussd.SetSession(sessionId, "amount", "1000")
func SetSession(sessionId, key, value string) {
	val, ok := session[sessionId]
	if !ok {
		val = make(map[string]string)
	}

	val[key] = value
	session[sessionId] = val
}

// RemoveSession deletes a session and all its data.
// Called automatically when a menu step returns "END".
//
// Example:
//
//	// Clean up after user completes flow
//	ussd.RemoveSession(sessionId)
func RemoveSession(sessionId string) {
	delete(session, sessionId)
}

// ClearSession removes all sessions.
// Useful for testing or administrative cleanup.
//
// Example:
//
//	// In test setup
//	ussd.ClearSession()
func ClearSession() {
	session = make(map[string]map[string]string)
}
