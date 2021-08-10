package cache

import (
	"math/rand"
	"testing"
)

func TestCacheMapBasic(t *testing.T) {
	capacity := 2
	cm := newLRUCache(capacity)
	cases := []struct {
		key   string
		value string
	}{
		{"aaa", "111"},
		{"bbb", "222"},
		{"ccc", "333"},
	}
	t.Log("Test basic")
	for _, cas := range cases {
		_, ok := cm.Get(cas.key)
		if ok {
			t.Error()
		}
		cm.Add(cas.key, cas.value)
		v, ok := cm.Get(cas.key)
		if !ok || v.(string) != cas.value {
			t.Errorf("%v%v", ok, v)
		}
	}
	t.Log("Test capcity")
	outdatedCases := cases[:len(cases)-capacity]
	for _, cas := range outdatedCases {
		_, ok := cm.Get(cas.key)
		if ok {
			t.Errorf("%v%v", ok, cas.value)
		}
	}
}

func BenchmarkCacheMap_R90(b *testing.B) {
	benchmarkCacheMap(b, 0.8)
}
func BenchmarkCacheMap_R80(b *testing.B) {
	benchmarkCacheMap(b, 0.8)
}
func BenchmarkCacheMap_R70(b *testing.B) {
	benchmarkCacheMap(b, 0.8)
}
func BenchmarkCacheMapParallel_R90(b *testing.B) {
	benchmarkCacheMapParallel(b, 0.8)
}
func BenchmarkCacheMapParallel_R80(b *testing.B) {
	benchmarkCacheMapParallel(b, 0.8)
}
func BenchmarkCacheMapParallel_R70(b *testing.B) {
	benchmarkCacheMapParallel(b, 0.8)
}

func benchmarkCacheMap(b *testing.B, readFreq float32) {
	capacity := 50
	cm := newLRUCache(capacity)
	k := 0
	for i := 0; i < b.N; i++ {
		k += (rand.Intn(5) - 2)
		rand.Float32()
		if rand.Float32() > readFreq {
			v := rand.Intn(100)
			cm.Add(k, v)
		} else {
			cm.Get(k)
		}
	}
}

func benchmarkCacheMapParallel(b *testing.B, readFreq float32) {
	capacity := 50
	cm := newLRUCache(capacity)
	b.RunParallel(func(p *testing.PB) {
		k := 0
		for p.Next() {
			k += (rand.Intn(5) - 2)
			rand.Float32()
			if rand.Float32() > readFreq {
				v := rand.Intn(100)
				cm.Add(k, v)
			} else {
				cm.Get(k)
			}
		}
	})
}
