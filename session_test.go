package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	s := NewSession("dummy id", 0)
	assert.NotNil(t, s)

	s.Set("abc", "abc")
	assert.Nil(t, s.Get("def"))
	assert.NotNil(t, s.Get("abc"))

	s.Set("struct", struct {
		name string
	}{
		"Jhon Doe",
	})
	_, err := s.GobEncode()
	assert.Error(t, err)

	s.Delete("struct")
	assert.Nil(t, s.Get("struct"))

	bs, err := s.GobEncode()
	assert.NoError(t, err)
	assert.NotNil(t, bs)
	assert.True(t, len(bs) > 15 /* time.Time */)

	s2 := NewSession("dummy id", 0)
	assert.NoError(t, s2.GobDecode(bs))
	assert.True(t, s.data.Expires.Equal(s2.data.Expires))
	assert.Equal(t, s.data.Values, s2.data.Values)
}
