package bitcask

import (
	"errors"
	"fmt"
	"sync"

	"github.com/0rz1/bitcask/cache"
	"github.com/0rz1/bitcask/set"
)

// communication
type comm struct {
	key   []byte
	value []byte
	loc   *location
	err   error
	res   chan<- *comm
}

type DB struct {
	cxt    *context
	cache  cache.Cache
	set    *set.Set
	app    *appender
	reader *reader
	loader *loader
	//comm
	readQ   chan *comm
	writeQ  chan *comm
	closeWG sync.WaitGroup
	closed  bool
}

func Open(path string, options ...Option) (*DB, error) {
	cxt := &context{path: path}
	if err := cxt.check(); err != nil {
		return nil, err
	}
	db := &DB{
		cxt:     cxt,
		set:     set.New(),
		app:     newAppender(cxt),
		reader:  newReader(cxt),
		loader:  newLoader(cxt),
		readQ:   make(chan *comm, 1),
		writeQ:  make(chan *comm, 1),
		closeWG: sync.WaitGroup{},
	}
	for _, opt := range options {
		if err := opt.custom(db); err != nil {
			return nil, err
		}
	}
	for _, opt := range []Option{defaultCacheOption, defaultLimitOption, defaultDiskOption} {
		if err := opt.custom(db); err != nil {
			if !errors.Is(err, ErrDuplicateOption) {
				return nil, err
			}
		}
	}
	if err := db.loader.load(db.set); err != nil {
		return nil, err
	}
	db.start()
	return db, nil
}

func (db *DB) Get(key string) (string, error) {
	if len(key) > db.cxt.limitOpt.MaxKeySize {
		return "", ErrKeyLenTooLong
	}
	if v, ok := db.cache.Get(key); ok {
		return v.(string), nil
	}
	if comp, ok := db.set.Get(key); ok {
		loc := comp.(*location)
		bs, err := db.read(loc)
		if err != nil {
			fmt.Println(*loc)
			fmt.Println(err)
			return "", err
		}
		v := string(bs)
		db.cache.Add(key, v)
		return v, nil
	}
	return "", ErrKeyNotFound
}

func (db *DB) Add(key, value string) error {
	if len(key) > db.cxt.limitOpt.MaxKeySize {
		return ErrKeyLenTooLong
	}
	if len(value) > db.cxt.limitOpt.MaxValueSize {
		return ErrValueLenTooLong
	}
	loc, err := db.write([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	db.set.Add(key, loc)
	db.cache.Add(key, value)
	return nil
}

func (db *DB) read(loc *location) ([]byte, error) {
	res := make(chan *comm)
	c := &comm{
		loc: loc,
		res: res,
	}
	db.readQ <- c
	r := <-res
	return r.value, r.err
}

func (db *DB) write(key, value []byte) (*location, error) {
	res := make(chan *comm)
	c := &comm{
		key:   key,
		value: value,
		res:   res,
	}
	db.writeQ <- c
	r := (<-res)
	return r.loc, r.err
}

func (db *DB) Close() {
	if !db.closed {
		close(db.readQ)
		close(db.writeQ)
		db.closeWG.Wait()
		db.reader.close()
		db.app.close()
		db.closed = true
	}
}

func (db *DB) start() {
	for i := 0; i < db.cxt.diskOpt.ReaderCnt; i++ {
		db.closeWG.Add(1)
		go func() {
			for c := range db.readQ {
				c.value, c.err = db.reader.read(c.loc)
				c.res <- c
			}
			db.closeWG.Done()
		}()
	}
	db.closeWG.Add(1)
	go func() {
		for c := range db.writeQ {
			c.loc, c.err = db.app.append(c.key, c.value)
			c.res <- c
		}
		db.closeWG.Done()
	}()
}
