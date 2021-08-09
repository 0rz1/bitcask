package db

import (
	"bytes"
	"testing"
)

func TestCacheMapBasic(t *testing.T) {
	capacity := 2
	cm := newCacheMap(capacity)
	cases := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("aaa"), []byte("111")},
		{[]byte("bbb"), []byte("222")},
		{[]byte("ccc"), []byte("333")},
	}
	t.Log("Test basic")
	for _, cas := range cases {
		_, ok := cm.get(cas.key)
		if ok {
			t.Error()
		}
		cm.put(cas.key, cas.value)
		v, ok := cm.get(cas.key)
		if !ok || !bytes.Equal(v, cas.value) {
			t.Errorf("%v%v", ok, v)
		}
	}
	t.Log("Test capcity")
	outdatedCases := cases[:len(cases)-capacity]
	for _, cas := range outdatedCases {
		_, ok := cm.get(cas.key)
		if ok {
			t.Errorf("%v%v", ok, cas.value)
		}
	}
}
