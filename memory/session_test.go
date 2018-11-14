package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	defer func() {
		assert.Nil(t, recover())
	}()
	s := NewSession("dummy")

	assert.Equal(t, "dummy", s.SessionID())
	assert.NoError(t, s.Set("name", "Jhon Doe"))

	assert.Equal(t, "Jhon Doe", s.Get("name"))
	assert.Nil(t, s.Get("pass"))

	s.Delete("name")
	assert.Nil(t, s.Get("name"))
	s.Delete("pass")
}
