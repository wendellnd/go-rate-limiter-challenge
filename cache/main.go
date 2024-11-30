package cache

type Cache interface {
	GetIncrement(key string) (int, error)
	SetExpire(key string, value int, duration int) error
	Increment(key string) error
}
