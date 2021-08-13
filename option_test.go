package bitcask

import (
	"testing"

	"github.com/0rz1/bitcask/cache"
)

func TestCacheOption(t *testing.T) {
	opt := &CacheOption{
		Capacity: 10,
	}
	db := &DB{cxt: &context{}}
	if err := opt.custom(db); err == nil {
		lruc := db.cache.(*cache.LRUCache)
		if lruc.Capacity() != opt.Capacity {
			t.Error()
		}
	} else {
		t.Error()
	}
}

func TestLimitOption(t *testing.T) {
	opt := &LimitOption{
		MaxFileSize:  10000,
		MaxKeySize:   10001,
		MaxValueSize: 10002,
	}
	db := &DB{cxt: &context{}}
	if err := opt.custom(db); err == nil {
		if db.cxt.limitOpt.MaxFileSize != opt.MaxFileSize {
			t.Error()
		} else if db.cxt.limitOpt.MaxKeySize != opt.MaxKeySize {
			t.Error()
		} else if db.cxt.limitOpt.MaxValueSize != opt.MaxValueSize {
			t.Error()
		}
	} else {
		t.Error()
	}
}

func TestDiskOption(t *testing.T) {
	opt := &DiskOption{
		ReaderCnt: 10,
		LoaderCnt: 1,
	}
	db := &DB{cxt: &context{}}
	if err := opt.custom(db); err == nil {
		if db.cxt.diskOpt.ReaderCnt != opt.ReaderCnt {
			t.Error()
		}
	} else {
		t.Error()
	}
}
