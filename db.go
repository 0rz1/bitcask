package bitcask

import (
	"github.com/0rz1/bitcask/cache"
	"github.com/0rz1/bitcask/set"
)

type DB struct {
	cxt   *context
	cache cache.Cache
	set   *set.Set
}

func New(path string, options ...Option) (*DB, error) {
	cxt := &context{path: path}
	db := &DB{
		cxt: cxt,
		set: set.New(),
	}
	for _, opt := range options {
		if err := opt.custom(db); err != nil {
			return nil, err
		}
	}
	if db.cache == nil {
		defaultCacheOption.custom(db)
	}
	if db.cxt.max_filesize == 0 {
		defaultLimitOption.custom(db)
	}
	return db, nil
}
