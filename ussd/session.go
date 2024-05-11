package ussd

var session map[string]map[string]string

func init() {
	session = make(map[string]map[string]string)
}

func HasSession(sessionId string) bool {
	_, ok := session[sessionId]
	return ok
}

func GetSession(sessionId string) map[string]string {
	if !HasSession(sessionId) {
		SetSession(sessionId, make(map[string]string))
	}

	return session[sessionId]
}

func SetSession(sessionId string, params map[string]string) {
	session[sessionId] = params
}

func RemoveSession(sessionId string) {
	delete(session, sessionId)
}

func ClearSession() {
	session = make(map[string]map[string]string)
}
