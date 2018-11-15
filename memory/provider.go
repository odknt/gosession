package memory

import (
	"fmt"
	"sync"

	"github.com/odknt/gosession"
)

// Provider is an implementation of session.Provider for in-memory.
type Provider struct {
	sessions map[string]*session.Session
	mutex    sync.RWMutex
}

// New returns a new Provider.
func New() *Provider {
	return &Provider{
		sessions: map[string]*session.Session{},
	}
}

// Init returns a new session, always returns error as nil.
func (p *Provider) Init(s *session.Session) error {
	p.sessions[s.ID()] = s
	return nil
}

// Read finds and returns a session by given session id.
// Returns a error if not found.
func (p *Provider) Read(sid string) (*session.Session, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	s, ok := p.sessions[sid]
	if !ok {
		return s, fmt.Errorf("not found session by given session id")
	}
	return s, nil
}

// Destroy removes a session by given session id.
func (p *Provider) Destroy(sid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.sessions[sid]; !ok {
		return fmt.Errorf("not found session by given session id")
	}

	delete(p.sessions, sid)
	return nil
}

// Commit nothing to do. returns nil always.
func (p *Provider) Commit(sid string) error {
	return nil
}

func init() {
	session.MustRegister("memory", New())
}
