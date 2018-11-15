package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/odknt/gosession"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	dir, err := ioutil.TempDir("", "gosession-test.")
	defer os.RemoveAll(dir)

	p := New(dir, "gosession-")
	assert.NotNil(t, p)

	// new session.
	s := session.NewSession("dummy", 0)
	assert.NoError(t, p.Init(s))

	// commit successful.
	s.Set("name", "Jhon Doe")
	assert.NoError(t, p.Commit("dummy"))

	// read session file.
	s, err = p.Read("dummy")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, "Jhon Doe", s.Get("name"))

	// destroy session.
	assert.NoError(t, p.Destroy("dummy"))

	// session file not found.
	s, err = p.Read("dummy")
	assert.Error(t, err)

	// invalid session file.
	assert.NoError(t, ioutil.WriteFile(p.getPath("invalid"), []byte("invalid"), os.FileMode(0777)))
	s, err = p.Read("invalid")
	assert.Error(t, err)

	// commit unknown session id.
	assert.Error(t, p.Commit("unknown"))

	// encode failed.
	s = session.NewSession("invalid", 0)
	assert.NoError(t, p.Init(s))
	s.Set("can't encode", struct{}{})
	assert.Error(t, p.Commit("invalid"))

	// destroy unknown session id.
	assert.Error(t, p.Destroy("unknown"))
	// destroy target file not found.
	s = session.NewSession("dummy", 0)
	assert.NoError(t, p.Init(s))
	assert.NoError(t, os.Remove(p.getPath("dummy")))
	assert.Error(t, p.Destroy("dummy"))
}
