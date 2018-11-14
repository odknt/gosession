package memory

import (
	"fmt"

	"github.com/odknt/gosession"
)

// Provider is an implementation of session.Provider for in-memory.
type Provider map[string]Session

// Init returns a new session, always returns error as nil.
func (p Provider) Init(sid string) (session.Session, error) {
	s := NewSession(sid)
	p[sid] = s
	return s, nil
}

// Read finds and returns a session by given session id.
// Returns a error if not found.
func (p Provider) Read(sid string) (session.Session, error) {
	s, ok := p[sid]
	if !ok {
		return s, fmt.Errorf("not found session by given session id")
	}
	return s, nil
}

// Destroy removes a session by given session id.
func (p Provider) Destroy(sid string) error {
	if _, ok := p[sid]; !ok {
		return fmt.Errorf("not found session by given session id")
	}
	delete(p, sid)
	return nil
}

func init() {
	session.MustRegister("memory", Provider{})
}
