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
	val, ok := session[sessionId]
	if !ok {
		return make(map[string]string)
	}

	return val
}

func SetSession(sessionId, key, value string) {
	val, ok := session[sessionId]
	if !ok {
		val = make(map[string]string)
	}

	val[key] = value
	session[sessionId] = val
}

func RemoveSession(sessionId string) {
	delete(session, sessionId)
}

func ClearSession() {
	session = make(map[string]map[string]string)
}
