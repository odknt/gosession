package memory

import (
	"testing"

	"github.com/odknt/gosession"
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
	p := New()

	for _, sid := range testSIDs {
		s := session.NewSession(sid, -1)
		err := p.Init(s)
		assert.NoError(t, err)
		assert.Equal(t, sid, s.ID())
	}

	for _, sid := range testSIDs {
		s, err := p.Read(sid)
		assert.NoError(t, err)
		assert.Equal(t, sid, s.ID())
	}
	_, err := p.Read("unknown session")
	assert.Error(t, err)

	for _, sid := range testSIDs {
		// Commit returns nil always.
		assert.Nil(t, p.Commit(sid))
	}

	for _, sid := range testSIDs {
		err := p.Destroy(sid)
		assert.NoError(t, err)
	}
	err = p.Destroy("unknown session")
	assert.Error(t, err)

	assert.Empty(t, p.sessions)
}
