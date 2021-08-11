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
func (a *appendFile) prepare() {
	if !a.freshed {
		var err error
		if a.file, err = uOpen(a.ft, a.no, a.cxt); err == nil {
			a.freshed = true
		}
	}
	if a.freshed && !a.aligned {
		if off, err := a.file.Seek(0, 2); err == nil {
			a.offset = int(off)
			a.aligned = true
		}
	}
}
func (a *appendFile) check(no int) {
	if a.prepared() {
		a.file.Close()
	}
	a.no = no
	a.freshed = false
	a.aligned = false
	a.prepare()
}
func (a *appendFile) write(bs []byte) bool {
	n, err := a.file.Write(bs)
	a.offset += n
	return err != nil
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

func (a *appender) prepare() {
	a.loc.prepare()
	a.dat.prepare()
}

func newAppender(cxt *context) *appender {
	a := &appender{
		cxt: cxt,
		loc: &appendFile{ft: FT_Location, cxt: cxt},
		dat: &appendFile{ft: FT_Data, cxt: cxt},
	}
	return a
}

func (a *appender) append(key, value []byte) {
	if a.prepare(); !a.prepared() {
		return
	}
	locSize := locationSeqSize(len(key))
	datSize := len(value)
	if a.loc.exLimit(locSize) || a.dat.exLimit(datSize) {
		a.no++
		a.loc.check(a.no)
		a.dat.check(a.no)
		if a.prepare(); !a.prepared() {
			return
		}
	}
	loc := &location{
		fileno: a.no,
		offset: a.dat.offset,
		length: datSize,
	}
	if !a.dat.write(value) {
		return
	}
	locbs := loc.makeSeqWithKey(key)
	if !a.loc.write(locbs) {
		return
	}
}
