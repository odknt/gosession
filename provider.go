package session

type Provider interface {
    Init(sid string) (Session, error)
    Read(sid string) (Session, error)
    Destroy(sid string) error
}
