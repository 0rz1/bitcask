package bitcask

import (
	"github.com/0rz1/bitcask/cache"
	"github.com/0rz1/bitcask/set"
)

type DB struct {
	cxt    *context
	cache  cache.Cache
	set    *set.Set
	app    *appender
	reader *reader
	loader *loader
}

func Open(path string, options ...Option) (*DB, error) {
	cxt := &context{path: path}
	if err := cxt.check(); err != nil {
		return nil, err
	}
	db := &DB{
		cxt:    cxt,
		set:    set.New(),
		app:    newAppender(cxt),
		reader: newReader(cxt),
		loader: newLoader(cxt),
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
	if err := db.loader.load(db.set); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() {

}

func (db *DB) GetSingle(key string) (string, error) {
	if v, ok := db.cache.Get(key); ok {
		return v.(string), nil
	}
	if comp, ok := db.set.Get(key); ok {
		loc := comp.(*location)
		bs := db.reader.read(loc)
		if len(bs) == 0 {
			return "", ErrDiskRD
		}
		v := string(bs)
		db.cache.Add(key, v)
		return v, nil
	}
	return "", nil
}

func (db *DB) AddSingle(key, value string) error {
	loc, err := db.app.append([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	db.set.Add(key, loc)
	db.cache.Add(key, value)
	return nil
}
