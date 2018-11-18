package session

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var providers = map[string]Provider{}

// Option is options when Manager creating session.
type Option struct {
	// defaults:
	//   Cookie:   "gosessionid"
	//   MaxAge:   0

	Cookie string
	MaxAge int
	SIDLen int
}

func setDefaults(opts Option) Option {
	if opts.Cookie == "" {
		opts.Cookie = "gosessionid"
	}
	if opts.SIDLen == 0 {
		opts.SIDLen = 32
	}
	return opts
}

// Manager controls session by using session provider.
type Manager struct {
	provider Provider
	opts     Option
}

// NewManager returns a new Manager given a provider name, cookie name, max lifetime.
func NewManager(providerName string, opts Option) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q", providerName)
	}
	return &Manager{
		provider: provider,
		opts:     setDefaults(opts),
	}, nil
}

// newSID returns a new session id.
func (m *Manager) newSID() (string, error) {
	sid := make([]byte, m.opts.SIDLen)
	if _, err := rand.Read(sid); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sid), nil
}

// newSession returns a new session.
func (m *Manager) newSession(w http.ResponseWriter) (*Session, error) {
	sid, err := m.newSID()
	if err != nil {
		return nil, err
	}

	s := NewSession(sid, m.opts.MaxAge)
	if err := m.provider.Init(s); err != nil {
		return nil, err
	}
	go time.AfterFunc(time.Duration(m.opts.MaxAge)*time.Second, func() { m.provider.Destroy(sid) })

	cookie := &http.Cookie{
		Name:     m.opts.Cookie,
		Value:    url.QueryEscape(sid),
		HttpOnly: true,
		MaxAge:   m.opts.MaxAge,
	}
	http.SetCookie(w, cookie)

	return s, nil
}

// Start finds and returns a session given a request cookie.
// If not found then returns a new session.
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// if the cookie hasn't session id then sets a new session.
	cookie, err := r.Cookie(m.opts.Cookie)
	if err != nil || cookie.Value == "" {
		return m.newSession(w)
	}

	// get session id.
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return m.newSession(w)
	}

	// find a session by the session id. if not found return new session.
	s, err := m.provider.Read(sid)
	if err != nil {
		return m.newSession(w)
	}

	// expires is before now.
	if s.data.Expires.Before(time.Now()) {
		m.Destroy(s)
		return m.newSession(w)
	}

	go time.AfterFunc(time.Until(s.data.Expires), func() {
		m.provider.Destroy(sid)
	})

	return s, nil
}

// Destroy removes a session from provider.
func (m *Manager) Destroy(s *Session) error {
	return m.provider.Destroy(s.ID())
}

// Commit calls Provider.Commit that session to be persistence.
func (m *Manager) Commit(s *Session) error {
	return m.provider.Commit(s.ID())
}

// Register registers a provider with specified name.
func Register(name string, provider Provider) error {
	if provider == nil {
		return errors.New("Register given provider is nil")
	}
	if _, dup := providers[name]; dup {
		return fmt.Errorf("Register called twice for provider: %s", name)
	}
	providers[name] = provider
	return nil
}

// MustRegister registers a provider with specified name.
// If register failed then panic.
func MustRegister(name string, provider Provider) {
	if err := Register(name, provider); err != nil {
		panic(err)
	}
}
