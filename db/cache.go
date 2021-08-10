package db

import (
	"container/list"
	"sync"
)

type cacheEntry struct {
	key   string
	value []byte
}

type cacheMap struct {
	capacity int
	recent   *list.List
	store    *sync.Map
	mux      *sync.Mutex
}

func newCacheMap(capacity int) *cacheMap {
	return &cacheMap{
		capacity: capacity,
		recent:   list.New(),
		store:    &sync.Map{},
		mux:      &sync.Mutex{},
	}
}

func (c *cacheMap) get(key string) ([]byte, bool) {
	iele, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	ele := iele.(*list.Element)
	c.mux.Lock()
	c.recent.MoveToFront(ele)
	ent := ele.Value.(*cacheEntry)
	c.mux.Unlock()
	return ent.value, true
}

func (c *cacheMap) add(key string, value []byte) {
	ent := &cacheEntry{
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
		c.store.Delete(ient.(*cacheEntry).key)
	}
}
