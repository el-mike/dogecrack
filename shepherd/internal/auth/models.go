package auth

const sessionIdCookie = "session_id"
const contextUserKey = "user"

// ContextUser - a Context value attached to authenticated requests.
type ContextUser struct {
	sessionId string
	userId    string
}
