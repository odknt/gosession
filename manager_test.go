package session

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New("invalid", "dummy", -1)
	assert.Error(t, err)
	m, err := New("memory", "dummy", -1)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestNewSession(t *testing.T) {
	w := httptest.NewRecorder()

	// returns error on newSession.
	m, _ := New("error", "dummy", -1)
	_, err := m.newSession(w)
	assert.Error(t, err)

	// success new session.
	m, _ = New("memory", "dummy", -1)
	s, err := m.newSession(w)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestStart(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	m, _ := New("memory", "dummy", 86400)

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
	assert.NotEqual(t, sid, invalid.SessionID())

	// session not found
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.AddCookie(&http.Cookie{Name: "dummy", Value: "unknown"})
	notfound, err := m.Start(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, notfound)
	assert.NotEqual(t, sid, notfound.SessionID())
	assert.NotEqual(t, invalid.SessionID(), notfound.SessionID())

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
	assert.Equal(t, sid, s.SessionID())
}

func TestReadAndDestroy(t *testing.T) {
	m, _ := New("memory", "dummy", 86400)
	s := newSession("invalid")
	assert.Error(t, m.Destroy(s))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	s, _ = m.Start(w, r)
	assert.NoError(t, m.Destroy(s))
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
