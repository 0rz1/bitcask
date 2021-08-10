package db

import (
	"os"
)

type diskWriter struct {
	cfg         *Config
	no          int
	locOff      int64
	valOff      int64
	locFile     *os.File
	valFile     *os.File
	valOutdated bool
	locOutdated bool
	offOutdated bool
}

func newDiskWriter(cfg *Config) *diskWriter {
	w := &diskWriter{
		cfg:         cfg,
		valOutdated: true,
		locOutdated: true,
		offOutdated: true,
	}
	return w
}

func (w *diskWriter) outdated() bool {
	return w.valOutdated || w.locOutdated || w.offOutdated
}

func (w *diskWriter) outdate() {
	w.valOutdated = true
	w.locOutdated = true
	w.offOutdated = true
}

func (w *diskWriter) write(key, value []byte, ch chan<- *result) {
	if w.refresh(); w.outdated() {
		return
	}
	seqSize := LocationSeqSize(len(key))
	if int(w.locOff)+seqSize > MaxFileSize ||
		int(w.valOff)+len(value) > MaxFileSize {
		if w.locOff == 0 && w.valOff == 0 {
			// TBD add error
			ch <- newResult(nil, nil)
			return
		}
		w.no++
		w.locOff = 0
		w.valOff = 0
		w.outdate()
		if w.refresh(); w.outdated() {
			// TBD add error
			ch <- newResult(nil, nil)
			return
		}
	}
	loc := &location{}
	loc.fileno = w.no
	loc.offset = int(w.valOff)
	loc.length = len(value)
	n, err := w.valFile.Write(value)
	w.valOff += int64(n)
	if err != nil {
		ch <- newResult(nil, err)
	}
	locbs := loc.MakeSeqWithKey(key)
	n, err = w.locFile.Write(locbs)
	w.locOff += int64(n)
	if err != nil {
		ch <- newResult(loc, err)
	}
}

func (w *diskWriter) refresh() {
	if !w.outdated() {
		return
	}
	var err error
	if w.locOutdated {
		if w.locFile != nil {
			w.locFile.Close()
			w.locFile = nil
		}
		w.locFile, err = w.cfg.OpenLocFile(w.no, false)
		if err != nil {
			w.locOutdated = false
		}
	}
	if w.valOutdated {
		if w.valFile != nil {
			w.valFile.Close()
			w.valFile = nil
		}
		w.valFile, err = w.cfg.OpenValFile(w.no, false)
		if err != nil {
			w.valOutdated = false
		}
	}

	if !w.valOutdated && !w.locOutdated {
		w.valOff, err = w.valFile.Seek(0, 2)
		if err != nil {
			return
		}
		w.locOff, err = w.locFile.Seek(0, 2)
		if err != nil {
			return
		}
		w.offOutdated = false
	}
}
