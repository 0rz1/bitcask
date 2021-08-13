package bitcask

import (
	"errors"

	"github.com/0rz1/bitcask/cache"
)

type Option interface {
	custom(db *DB) error
}

var _ Option = &CacheOption{}
var _ Option = &LimitOption{}
var _ Option = &DiskOption{}

type CacheOption struct {
	Capacity int
}

func (opt *CacheOption) custom(db *DB) error {
	if db.cache != nil {
		return ErrDuplicateOption
	}
	if opt.Capacity < 5 {
		return errors.New("capacity less than 5")
	}
	db.cache = cache.NewLRUCache(opt.Capacity)
	return nil
}

type LimitOption struct {
	MaxFileSize  int
	MaxKeySize   int
	MaxValueSize int
}

func (opt *LimitOption) custom(db *DB) error {
	if db.cxt.limitOpt.MaxFileSize == 0 {
		db.cxt.limitOpt = *opt
	} else {
		return ErrDuplicateOption
	}
	return nil
}

type DiskOption struct {
	ReaderCnt int
	LoaderCnt int
}

func (opt *DiskOption) custom(db *DB) error {
	if db.cxt.diskOpt.ReaderCnt == 0 {
		db.cxt.diskOpt = *opt
		return nil
	} else {
		return ErrDuplicateOption
	}
}
