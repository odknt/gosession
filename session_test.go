package session

type inMemorySession struct {
	sid    string
	values map[string]interface{}
}

func newSession(sid string) Session {
	return inMemorySession{
		sid:    sid,
		values: map[string]interface{}{},
	}
}

func (s inMemorySession) Set(key string, value interface{}) error {
	s.values[key] = value
	return nil
}

func (s inMemorySession) Get(key string) interface{} {
	return s.values[key]
}

func (s inMemorySession) Delete(key string) {
	delete(s.values, key)
}

func (s inMemorySession) SessionID() string {
	return s.sid
}
