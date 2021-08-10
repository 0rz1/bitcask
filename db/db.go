package db

type DB struct {
	cfg     *Config
	reader  *diskReader
	writer  *diskWriter
	cache   *cacheMap
	locates *locationMap
	quit    chan int
}

func Open(cfg *Config) (*DB, error) {
	if !cfg.Check() {
		//TBD add error
		return nil, nil
	}
	db := &DB{
		cfg:     cfg,
		reader:  newDiskReader(cfg),
		writer:  newDiskWriter(cfg),
		cache:   newCacheMap(100),
		locates: newLocationMap(),
		quit:    make(chan int),
	}
	// go db.loop()
	return db, nil
}

func (db *DB) Close() {
	// db.quit <- 1
}

func (db *DB) Read(key string) (string, error) {
	v, ok := db.cache.get(key)
	if ok {
		// Get In Cache
		return string(v), nil
	}
	loc, ok := db.locates.get(key)
	if !ok {
		// No Key
		return "", nil
	}
	resCh := make(chan *result)
	go db.reader.read(loc, resCh)
	res := <-resCh
	if res.err == nil {
		v := string(res.value.([]byte))
		return v, res.err
	} else {
		return "", nil
	}
}

func (db *DB) Write(key, value string) error {
	resCh := make(chan *result)
	go db.writer.write([]byte(key), []byte(value), resCh)
	res := <-resCh
	if res.err != nil {
		return res.err
	}
	loc := res.value.(*location)
	db.locates.add(key, loc)
	db.cache.add(key, []byte(value))
	return nil
}

func (db *DB) loop() {
	// for {
	// 	select {
	// 	case <-db.quit:
	// 		break
	// 	}
	// }
}
