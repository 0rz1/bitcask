package bitcask

type reader struct {
	cxt *context
}

func newReader(cxt *context) *reader {
	return &reader{
		cxt: cxt,
	}
}

func (r *reader) read(loc *location) ([]byte, error) {
	f, err := uOpen(FT_Data, loc.fileno, r.cxt)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bs := make([]byte, loc.length)
	n, err := f.ReadAt(bs, int64(loc.offset))
	if n != loc.length {
		return nil, err
	} else {
		return bs, nil
	}
}

func (r *reader) close() {

}
