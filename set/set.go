package set

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Comparable interface {
	Compare(other Comparable) int
}

type entry struct {
	ptr unsafe.Pointer // Comparable
}

func makeEntry(val Comparable) *entry {
	return &entry{
		ptr: unsafe.Pointer(&val),
	}
}

type Key interface{}

type Set struct {
	store *sync.Map
}

func New() *Set {
	return &Set{
		store: &sync.Map{},
	}
}

func (set *Set) Add(key Key, val Comparable) {
	ent := makeEntry(val)
	actual, loaded := set.store.LoadOrStore(key, ent)
	if loaded {
		e := actual.(*entry)
		for {
			p := e.ptr
			v := *(*Comparable)(p)
			if val.Compare(v) <= 0 {
				return
			}
			if atomic.CompareAndSwapPointer(&e.ptr, p, ent.ptr) {
				return
			}
		}
	}
}

func (set *Set) Get(key Key) (Comparable, bool) {
	_e, ok := set.store.Load(key)
	if !ok {
		return nil, false
	}
	e := _e.(*entry)
	v := *(*Comparable)(e.ptr)
	return v, ok
}
