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

func Open(path string, options ...Option) (*DB, error) {
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

func (db *DB) Close() {

}

func (db *DB) Get(key string) (value string, err error) {
	res := asyncResponse(key, func(r *request) { db.get(r) })
	if res.err != nil {
		return "", res.err
	}
	return res.ret.(string), nil
}

func (db *DB) Add(key, value string) (err error) {
	res := asyncResponse(&struct {
		key   string
		value string
	}{key, value}, func(r *request) { db.add(r) })
	if res.err != nil {
		return res.err
	}
	return nil
}

func (db *DB) get(req *request) {
	// key := req.param.(string)
}

func (db *DB) add(req *request) {
}
