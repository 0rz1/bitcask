package db

// prev <-> cur <-> next
type storeNode struct {
	key   string
	value string
	prev  *storeNode
	next  *storeNode
}

func newStoreNode(key, value string) *storeNode {
	sn := &storeNode{
		key:   key,
		value: value,
	}
	sn.prev = sn
	sn.next = sn
	return sn
}

func storeDel(sn *storeNode) {
	sn.prev.next = sn.next
	sn.next.prev = sn.prev
}

func storeIns(cur, sn *storeNode) *storeNode {
	if cur == nil {
		return sn
	}
	if sn == nil {
		return cur
	}
	sn.next = cur.next
	sn.prev = cur
	cur.next = sn
	return sn
}

type cacheMap struct {
	capacity int
	size     int
	store    *storeNode
	locate   map[string]*storeNode
}

func newCacheMap(capacity int) *cacheMap {
	var cm = cacheMap{
		capacity: capacity,
		size:     0,
		locate:   make(map[string]*storeNode, capacity),
		store:    nil,
	}
	return &cm
}

func (c *cacheMap) get(key []byte) ([]byte, bool) {
	sn, ok := c.locate[string(key)]
	if !ok {
		return nil, ok
	}
	// if sn != c.store {
	// 	storeDel(sn)
	// 	c.store = storeIns(c.store, sn)
	// }
	return []byte(sn.value), ok
}

func (c *cacheMap) put(key, value []byte) {
	skey := string(key)
	svalue := string(value)
	sn, ok := c.locate[skey]
	if ok {
		if sn == c.store {
			return
		}
		storeDel(sn)
		sn.value = svalue
	} else {
		c.size++
		if c.size > c.capacity {
			nsn := c.store.next
			storeDel(nsn)
			delete(c.locate, nsn.key)
			c.size--
		}
		sn = newStoreNode(skey, svalue)
	}
	c.store = storeIns(c.store, sn)
	c.locate[skey] = sn
}
