package session

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Session represents a session.
type Session struct {
	sid  string
	data sessionData
}

// NewSession returns a new Session.
func NewSession(sid string, maxAge int) *Session {
	expires := time.Now().Add(time.Duration(maxAge) * time.Second)
	return &Session{
		sid: sid,
		data: sessionData{
			Expires: expires,
			Values:  map[string]interface{}{},
		},
	}
}

// Set sets a value with specified key and always returns nil.
func (s Session) Set(key string, value interface{}) {
	s.data.Values[key] = value
}

// Get gets a value by given key.
func (s Session) Get(key string) interface{} {
	return s.data.Values[key]
}

// Delete removes a value by given key.
func (s Session) Delete(key string) {
	delete(s.data.Values, key)
}

// ID returns myself session id.
func (s Session) ID() string {
	return s.sid
}

// Expired returns true if the session is expired.
func (s Session) Expired() bool {
	// expires is before now.
	return s.data.Expires.Before(time.Now())
}

// GobEncode implements the gob.GobEncoder interface.
func (s Session) GobEncode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(s.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements the gob.GobDecoder interface.
func (s *Session) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	return dec.Decode(&s.data)
}

type sessionData struct {
	Expires time.Time
	Values  map[string]interface{}
}

func init() {
	gob.Register(sessionData{})
}
