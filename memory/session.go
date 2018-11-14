package memory

// Session represents a session.
type Session struct {
	sid    string
	values map[string]interface{}
}

// NewSession returns a new session.
func NewSession(sid string) Session {
	return Session{
		sid:    sid,
		values: map[string]interface{}{},
	}
}

// Set sets a value with specified key and always returns nil.
func (s Session) Set(key string, value interface{}) error {
	s.values[key] = value
	return nil
}

// Get gets a value by given key.
func (s Session) Get(key string) interface{} {
	return s.values[key]
}

// Delete removes a value by given key.
func (s Session) Delete(key string) {
	delete(s.values, key)
}

// SessionID returns myself session id.
func (s Session) SessionID() string {
	return s.sid
}
