package session

type Session interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string)
	SessionID() string
}
