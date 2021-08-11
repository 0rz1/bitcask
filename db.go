package bitcask

import (
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
	diskOpt DiskOption
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
		readQ:  make(chan *comm),
		writeQ: make(chan *comm),
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
	if db.diskOpt.readerCnt == 0 {
		defaultDiskOption.custom(db)
	}
	if err := db.loader.load(db.set, db.diskOpt.loaderCnt); err != nil {
		return nil, err
	}
	db.start()
	return db, nil
}

func (db *DB) GetSingle(key string) (string, error) {
	if v, ok := db.cache.Get(key); ok {
		return v.(string), nil
	}
	if comp, ok := db.set.Get(key); ok {
		loc := comp.(*location)
		bs := db.read(loc)
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
	loc, err := db.write([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	db.set.Add(key, loc)
	db.cache.Add(key, value)
	return nil
}

func (db *DB) read(loc *location) []byte {
	res := make(chan *comm)
	c := &comm{
		loc: loc,
		res: res,
	}
	db.readQ <- c
	return (<-res).value
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
	close(db.readQ)
	close(db.writeQ)
}

func (db *DB) start() {
	for i := 0; i < db.diskOpt.readerCnt; i++ {
		go func() {
			for c := range db.readQ {
				c.value = db.reader.read(c.loc)
				c.err = nil
				c.res <- c
			}
		}()
	}
	go func() {
		for c := range db.writeQ {
			c.loc, c.err = db.app.append(c.key, c.value)
			c.res <- c
		}
	}()
}
