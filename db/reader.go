package db

type diskReader struct {
	cfg *Config
}

func (r *diskReader) readLocation(loc *location,
	ch chan<- *result) {
	f, err := r.cfg.OpenLocFile(loc.fileno, true)
	if err != nil {
		ch <- newResult(nil, err)
		return
	}
	bs := make([]byte, loc.length)
	_, err = f.ReadAt(bs, int64(loc.offset))
	if err != nil {
		ch <- newResult(nil, err)
	} else {
		ch <- newResult(bs, err)
	}
}
