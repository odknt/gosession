package session

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	opts := Option{Cookie: "dummy", MaxAge: -1}
	_, err := NewManager("invalid", opts)
	assert.Error(t, err)
	m, err := NewManager("memory", opts)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestManagerNewSession(t *testing.T) {
	w := httptest.NewRecorder()
	opts := Option{MaxAge: -1}

	// returns error on newSession.
	m, _ := NewManager("error", opts)
	_, err := m.newSession(w)
	assert.Error(t, err)

	// success new session.
	m, _ = NewManager("memory", opts)
	s, err := m.newSession(w)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestStart(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()

	opts := Option{Cookie: "dummy", MaxAge: 86400}
	m, _ := NewManager("memory", opts)

	// cookie not found.
	s, err := m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// get session id from response.
	sid := getSessionID(t, w)

	// cookie value unescape failed.
	w = httptest.NewRecorder()
	r.AddCookie(&http.Cookie{Name: "dummy", Value: "%2"})
	invalid, err := m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, invalid)
	assert.NotEqual(t, sid, invalid.ID())

	// session not found
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.AddCookie(&http.Cookie{Name: "dummy", Value: "unknown"})
	notfound, err := m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, notfound)
	assert.NotEqual(t, sid, notfound.ID())
	assert.NotEqual(t, invalid.ID(), notfound.ID())

	// exists session.
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.AddCookie(&http.Cookie{
		Name:     "dummy",
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})
	s, err = m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, sid, s.ID())

	// returns new session instead of expired session.
	s.data.Expires, _ = time.Parse("2006-01-02", "2000-01-01")
	m.provider.Init(s)
	olds, err := m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, olds)
	assert.NotEqual(t, s.ID(), olds.ID())

	// auto destroy after expires.
	s.data.Expires = time.Now().Add(10 * time.Millisecond)
	m.provider.Init(s)
	s, _ = m.Start(w, r)
	time.Sleep(11 * time.Millisecond)
	olds, err = m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, olds)
	assert.NotEqual(t, s.ID(), olds.ID())
}

func TestDestroy(t *testing.T) {
	opts := Option{Cookie: "dummy", MaxAge: 86400}
	m, _ := NewManager("memory", opts)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	s, _ := m.Start(w, r)
	assert.NoError(t, m.Destroy(s))
}

func TestCommit(t *testing.T) {
	opts := Option{Cookie: "dummy", MaxAge: 86400}
	m, _ := NewManager("memory", opts)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	s, _ := m.Start(w, r)
	assert.NoError(t, m.Commit(s))
}

func TestRegister(t *testing.T) {
	// must sets nil.
	assert.NotNil(t, Register("memory", nil))
	// must call only by provider name.
	assert.NotNil(t, Register("memory", inMemoryProvider{}))
}

func TestMustRegister(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	MustRegister("memory", inMemoryProvider{})
}

func getSessionID(t *testing.T, w *httptest.ResponseRecorder) string {
	var sid string
	if header, ok := w.HeaderMap["Set-Cookie"]; ok {
		r := &http.Request{Header: http.Header{"Cookie": strings.Split(header[0], "; ")}}
		// get cookie from response.
		cookie, err := r.Cookie("dummy")
		assert.NoError(t, err)
		assert.Equal(t, "dummy", cookie.Name)
		assert.NotEmpty(t, cookie.Value)

		// get session id from response.
		sid, err = url.QueryUnescape(cookie.Value)
		assert.NoError(t, err)
		assert.NotEmpty(t, sid)
	}

	return sid
}

func init() {
	MustRegister("memory", inMemoryProvider{})
	MustRegister("error", errorProvider{})
}
