package session

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/rs/xid"
)

var providers = map[string]Provider{}

// Manager controls session by using session provider.
type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
}

// New returns a new Manager given a provider name, cookie name, max lifetime.
func New(providerName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q", providerName)
	}
	return &Manager{
		provider:    provider,
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
	}, nil
}

// newSID returns a new session id.
func (m *Manager) newSID() string {
	guid := xid.New()
	return guid.String()
}

// newSession returns a new session.
func (m *Manager) newSession(w http.ResponseWriter) (Session, error) {
	sid := m.newSID()
	s, err := m.provider.Init(sid)
	if err != nil {
		return nil, err
	}
	go time.AfterFunc(time.Duration(m.maxLifeTime)*time.Second, func() { m.provider.Destroy(sid) })

	cookie := &http.Cookie{
		Name:     m.cookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(m.maxLifeTime),
	}
	http.SetCookie(w, cookie)

	return s, nil
}

// Start finds and returns a session given a request cookie.
// If not found then returns a new session.
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) (Session, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// if the cookie hasn't session id then sets a new session.
	cookie, err := r.Cookie(m.cookieName)
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

	return s, nil
}

// Destroy removes a session from provider.
func (m *Manager) Destroy(s Session) error {
	return m.provider.Destroy(s.SessionID())
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
