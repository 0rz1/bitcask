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
	if db.cxt.max_filesize == 0 {
		db.cxt.max_filesize = opt.MaxFileSize
		if db.cxt.max_filesize < 1000 {
			return errors.New("filesize less than 1000")
		}
	} else {
		return ErrDuplicateOption
	}
	if db.cxt.max_keysize == 0 {
		db.cxt.max_keysize = opt.MaxKeySize
		if db.cxt.max_keysize < 10 {
			return errors.New("keysize less than 10")
		}
	} else {
		return ErrDuplicateOption
	}
	if db.cxt.max_valuesize == 0 {
		db.cxt.max_valuesize = opt.MaxValueSize
		if db.cxt.max_valuesize < 100 {
			return errors.New("valuesize less than 10")
		}
	} else {
		return ErrDuplicateOption
	}
	return nil
}
