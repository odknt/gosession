package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	session "github.com/odknt/gosession"
	"github.com/odknt/gosession/memory"
)

// Provider is an implementation of session.Provider for file.
type Provider struct {
	parent *memory.Provider

	dir    string
	prefix string
	mutex  sync.RWMutex
}

// New returns a new Provider for file system.
// Writes session to a file with given prefix in given directory.
func New(dir, prefix string) *Provider {
	return &Provider{
		parent: memory.New(),
		dir:    dir,
		prefix: prefix,
	}
}

// Init returns a new session.
func (p *Provider) Init(s *session.Session) error {
	p.parent.Init(s)
	return p.save(s)
}

// Read finds and returns a session by given session id.
// Returns a error if not found.
func (p *Provider) Read(sid string) (*session.Session, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	bs, err := ioutil.ReadFile(p.getPath(sid))
	if err != nil {
		return nil, err
	}

	s := session.NewSession(sid, 0)
	if err := s.GobDecode(bs); err != nil {
		return nil, err
	}

	p.parent.Init(s)
	return s, nil
}

// Destroy removes a session by given session id.
// No-writes changes to a file until call Commit.
func (p *Provider) Destroy(sid string) error {
	_, err := p.parent.Read(sid)
	if err != nil {
		return err
	}
	if err := os.Remove(p.getPath(sid)); err != nil {
		return err
	}
	return p.parent.Destroy(sid)
}

// Commit writes session to a file.
func (p *Provider) Commit(sid string) error {
	s, err := p.parent.Read(sid)
	if err != nil {
		return err
	}

	return p.save(s)
}

// Cleanup cleans all sessions.
func (p *Provider) Cleanup() error {
	if err := filepath.Walk(p.dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(p.dir, fpath)
		if rel == "." {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		base := filepath.Base(fpath)
		if strings.HasPrefix(base, p.prefix) {
			sid := strings.TrimPrefix(base, p.prefix)
			s, err := p.Read(sid)
			if err == nil && s.Expired() {
				p.Destroy(sid)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (p *Provider) save(s *session.Session) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	bs, err := s.GobEncode()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.getPath(s.ID()), bs, os.FileMode(0600))
}

func (p *Provider) getPath(sid string) string {
	return filepath.Join(p.dir, p.prefix+sid)
}
