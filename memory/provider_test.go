package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSIDs = []string{
	"",
	"dummy",
	"has space value",
	"query%20escaped%20string",
	"invalid%2escaped%2string",
}

func TestProvider(t *testing.T) {
	p := Provider{}

	for _, sid := range testSIDs {
		s, err := p.Init(sid)
		assert.NoError(t, err)
		assert.Equal(t, sid, s.SessionID())
	}

	for _, sid := range testSIDs {
		s, err := p.Read(sid)
		assert.NoError(t, err)
		assert.Equal(t, sid, s.SessionID())
	}
	_, err := p.Read("unknown session")
	assert.Error(t, err)

	for _, sid := range testSIDs {
		err := p.Destroy(sid)
		assert.NoError(t, err)
	}
	err = p.Destroy("unknown session")
	assert.Error(t, err)

	assert.Empty(t, p)
}
