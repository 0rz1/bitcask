package bitcask

type reader struct {
	cxt *context
}

func newReader(cxt *context) *reader {
	return &reader{
		cxt: cxt,
	}
}

func (r *reader) read(loc *location) {
	f, err := uOpen(FT_Location, loc.fileno, r.cxt)
	if err != nil {
		return
	}
	bs := make([]byte, loc.length)
	_, err = f.ReadAt(bs, int64(loc.offset))
	if err != nil {
		return
	}
}