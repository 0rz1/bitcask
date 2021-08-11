package bitcask

import (
	"bytes"
	"testing"
)

func TestLocCompare(t *testing.T) {
	loc0 := &location{
		fileno: 0,
		offset: 100,
		length: 10,
	}
	loc1 := &location{
		fileno: 1,
		offset: 10,
		length: 1,
	}
	loc2 := &location{
		fileno: 1,
		offset: 11,
		length: 1,
	}
	if loc0.Compare(loc1) >= 0 {
		t.Error()
	}
	if loc1.Compare(loc2) >= 0 {
		t.Error()
	}
	if loc2.Compare(loc2) != 0 {
		t.Error()
	}
}

func TestLocSeq(t *testing.T) {
	loc := &location{
		fileno: 0,
		offset: 100,
		length: 10,
	}
	key := []byte("abc")
	bs := loc.makeSeqWithKey(key)
	if l, k, ok := makeLocationAndKey(bs); ok {
		if !bytes.Equal(key, k) || loc.Compare(l) != 0 {
			t.Error()
		}
	} else {
		t.Error()
	}
}

func TestLocSeq1(t *testing.T) {
	loc := &location{
		fileno: 0,
		offset: 100,
		length: 10,
	}
	key := []byte("abc")
	bs := loc.makeSeqWithKey(key)
	bs[len(bs)-1] = 1
	if _, _, ok := makeLocationAndKey(bs); ok {
		t.Error()
	}
}
