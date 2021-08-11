package bitcask

import (
	"testing"

	"github.com/0rz1/bitcask/cache"
)

func TestCacheOption(t *testing.T) {
	opt := &CacheOption{
		Capacity: 10,
	}
	if db, err := NewDB("", opt); err == nil {
		lruc := db.cache.(*cache.LRUCache)
		if lruc.Capacity() != opt.Capacity {
			t.Error()
		}
	} else {
		t.Error()
	}
	if db, err := NewDB(""); err == nil {
		lruc := db.cache.(*cache.LRUCache)
		if lruc.Capacity() != defaultCacheOption.Capacity {
			t.Error()
		}
	} else {
		t.Error()
	}
	if _, err := NewDB("", opt, opt); err != ErrDuplicateOption {
		t.Error()
	}
}

func TestLimitOption(t *testing.T) {
	opt := &LimitOption{
		MaxFileSize:  10000,
		MaxKeySize:   10001,
		MaxValueSize: 10002,
	}
	if db, err := NewDB("", opt); err == nil {
		if db.cxt.max_filesize != opt.MaxFileSize {
			t.Error()
		} else if db.cxt.max_keysize != opt.MaxKeySize {
			t.Error()
		} else if db.cxt.max_valuesize != opt.MaxValueSize {
			t.Error()
		}
	} else {
		t.Error()
	}
	if db, err := NewDB(""); err == nil {
		if db.cxt.max_filesize != defaultLimitOption.MaxFileSize {
			t.Error()
		} else if db.cxt.max_keysize != defaultLimitOption.MaxKeySize {
			t.Error()
		} else if db.cxt.max_valuesize != defaultLimitOption.MaxValueSize {
			t.Error()
		}
	} else {
		t.Error()
	}
	if _, err := NewDB("", opt, opt); err != ErrDuplicateOption {
		t.Error()
	}
}

func TestOptions(t *testing.T) {
	opt1 := &CacheOption{
		Capacity: 10,
	}
	opt2 := &LimitOption{
		MaxFileSize:  10000,
		MaxKeySize:   10001,
		MaxValueSize: 10002,
	}
	if _, err := NewDB("", opt1, opt2); err != nil {
		t.Error()
	}
	if _, err := NewDB(""); err != nil {
		t.Error()
	}
}
