package bitcask

import (
	"os"
)

type appendFile struct {
	ft      FileType
	no      int
	offset  int
	aligned bool
	freshed bool
	cxt     *context
	file    *os.File
}

func (a *appendFile) prepared() bool {
	return a.freshed && a.aligned
}
func (a *appendFile) prepare() error {
	if !a.freshed {
		var err error
		if a.file, err = uOpenAppend(a.ft, a.no, a.cxt); err == nil {
			a.freshed = true
		} else {
			return err
		}
	}
	if a.freshed && !a.aligned {
		if off, err := a.file.Seek(0, 2); err == nil {
			a.offset = int(off)
			a.aligned = true
		} else {
			return err
		}
	}
	return nil
}
func (a *appendFile) cut(no int) {
	if a.prepared() {
		a.file.Close()
	}
	a.no = no
	a.freshed = false
	a.aligned = false
	a.prepare()
}
func (a *appendFile) write(bs []byte) error {
	n, err := a.file.Write(bs)
	a.offset += n
	return err
}
func (a *appendFile) exLimit(sz int) bool {
	return a.offset+sz > a.cxt.max_filesize
}

type appender struct {
	cxt *context
	no  int
	loc *appendFile
	dat *appendFile
}

func (a *appender) prepared() bool {
	return a.loc.prepared() && a.dat.prepared()
}

func (a *appender) prepare() error {
	if err := a.loc.prepare(); err != nil {
		return err
	}
	return a.dat.prepare()
}

func newAppender(cxt *context) *appender {
	a := &appender{
		cxt: cxt,
		no:  cxt.maxno(),
		loc: &appendFile{ft: FT_Location, cxt: cxt},
		dat: &appendFile{ft: FT_Data, cxt: cxt},
	}
	return a
}

func (a *appender) append(key, value []byte) (*location, error) {
	if !a.prepared() {
		err := a.prepare()
		if err != nil {
			return nil, &Err{ErrNotReady, err}
		}
	}
	locSize := locationSeqSize(len(key))
	datSize := len(value)
	if a.loc.exLimit(locSize) || a.dat.exLimit(datSize) {
		a.no++
		a.loc.cut(a.no)
		a.dat.cut(a.no)
		err := a.prepare()
		if err != nil {
			return nil, &Err{ErrNotReady, err}
		}
	}
	loc := &location{
		fileno: a.no,
		offset: a.dat.offset,
		length: datSize,
	}
	if err := a.dat.write(value); err != nil {
		return nil, &Err{ErrWrite, err}
	}
	locbs := loc.makeSeqWithKey(key)
	if err := a.loc.write(locbs); err != nil {
		return nil, &Err{ErrWrite, err}
	}
	return loc, nil
}

func (a *appender) close() {
	if a.loc.file != nil {
		a.loc.file.Close()
	}
	if a.dat.file != nil {
		a.dat.file.Close()
	}
}
