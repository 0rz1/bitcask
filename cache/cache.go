package cache

type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

type Cache interface {
	Get(Key) (interface{}, bool)
	Add(Key, interface{})
}
