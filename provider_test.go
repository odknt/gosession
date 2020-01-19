package session

import (
	"errors"
	"fmt"
)

type inMemoryProvider map[string]*Session

func (p inMemoryProvider) Init(s *Session) error {
	p[s.ID()] = s
	return nil
}

func (p inMemoryProvider) Read(sid string) (*Session, error) {
	s, ok := p[sid]
	if !ok {
		return s, fmt.Errorf("not found session by given session id")
	}
	return s, nil
}

func (p inMemoryProvider) Destroy(sid string) error {
	if _, ok := p[sid]; !ok {
		return fmt.Errorf("not found session by given session id")
	}
	delete(p, sid)
	return nil
}

func (p inMemoryProvider) Commit(sid string) error {
	return nil
}

func (p inMemoryProvider) Cleanup() error {
	return nil
}

type errorProvider struct {
	inMemoryProvider
}

func (errorProvider) Init(*Session) error {
	return errors.New("initialize session failed")
}
