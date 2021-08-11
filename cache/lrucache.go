package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	capacity int
	recent   *list.List
	store    *sync.Map
	mux      *sync.Mutex
}

var _ Cache = &LRUCache{}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		recent:   list.New(),
		store:    &sync.Map{},
		mux:      &sync.Mutex{},
	}
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) Get(key Key) (interface{}, bool) {
	iele, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	ele := iele.(*list.Element)
	c.mux.Lock()
	c.recent.MoveToFront(ele)
	ent := ele.Value.(*entry)
	c.mux.Unlock()
	return ent.value, true
}

func (c *LRUCache) Add(key Key, value interface{}) {
	ent := &entry{
		key:   key,
		value: value,
	}
	iele, ok := c.store.Load(key)
	if ok {
		ele := iele.(*list.Element)
		c.mux.Lock()
		c.recent.MoveToFront(ele)
		ele.Value = ent
		c.mux.Unlock()
	} else {
		c.mux.Lock()
		ele := c.recent.PushFront(ent)
		c.mux.Unlock()
		c.store.Store(key, ele)
	}
	if c.recent.Len() > c.capacity {
		c.mux.Lock()
		ient := c.recent.Remove(c.recent.Back())
		c.mux.Unlock()
		c.store.Delete(ient.(*entry).key)
	}
}
