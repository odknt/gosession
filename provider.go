package session

// Provider is the interface for session providers.
type Provider interface {
	Init(session *Session) error
	Read(sid string) (*Session, error)
	Destroy(sid string) error
	Commit(sid string) error
}
