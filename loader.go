package bitcask

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/0rz1/bitcask/set"
)

type loader struct {
	cxt *context
}

func newLoader(cxt *context) *loader {
	return &loader{
		cxt: cxt,
	}
}

func (l *loader) load(set *set.Set, loaderCnt int) error {
	fq := make(chan *os.File)
	defer close(fq)
	eq := make(chan error)
	for i := 0; i < loaderCnt; i++ {
		go func() {
			for f := range fq {
				eq <- loadFile(f, set, l.cxt)
			}
		}()
	}
	for _, no := range l.cxt.filenos {
		if f, err := uOpen(FT_Location, no, l.cxt); err != nil {
			return err
		} else {
			fq <- f
		}
	}
	for e := range eq {
		if e != nil {
			return e
		}
	}
	return nil
}

func loadFile(f *os.File, set *set.Set, cxt *context) (err error) {
	defer f.Close()

	hlen := len(locSeqHeader)
	bufferSize := 4096
	maxLocSeqSize := locationSeqSize(cxt.max_keysize)
	buffer := make([]byte, maxLocSeqSize+bufferSize)

	var off = maxLocSeqSize
	var length = 0
	var ksz = 0
	var freshed = false
	var checkFresh func(int) bool = func(more int) bool {
		if off+more >= bufferSize {
			off -= bufferSize
			copy(buffer[off:], buffer[off+bufferSize:])
			freshed = false
			return false
		}
		return true
	}
	for n := 0; ; {
		if !freshed || off == n+maxLocSeqSize {
			n, err = f.Read(buffer[maxLocSeqSize:])
			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return err
			}
			freshed = true
			continue
		}
		switch length {
		case 0:
			if buffer[off] == locSeqHeader[0] {
				length = 1
			} else {
				off++
			}
		case 1:
			if !checkFresh(hlen) {
				continue
			}
			if bytes.Equal(buffer[off:off+hlen], locSeqHeader) {
				length = hlen
			} else {
				length = 0
				off++
			}
		case hlen:
			if !checkFresh(hlen + 16) {
				continue
			}
			length = hlen + 16
			ksz = int(binary.BigEndian.Uint32(
				buffer[off+length-4 : off+length]))
		case hlen + 16:
			if !checkFresh(hlen + 16 + ksz + 4) {
				continue
			}
			length = hlen + 16 + ksz + 4
			loc, key, ok := makeLocationAndKey(
				buffer[off : off+length])
			if ok {
				set.Add(string(key), loc)
				off += length
			} else {
				off++
			}
			length = 0
		default:
			panic("unknown error")
		}
	}
}
