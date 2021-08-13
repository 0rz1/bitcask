package bitcask

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"

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

func (l *loader) load(set *set.Set) (err error) {
	wg := sync.WaitGroup{}
	fq := make(chan *os.File)
	defer close(fq)
	for i := 0; i < l.cxt.diskOpt.LoaderCnt; i++ {
		go func() {
			for f := range fq {
				if err == nil {
					e := loadFile(f, set, l.cxt)
					if e != nil {
						err = e
					}
				}
				wg.Done()
			}
		}()
	}
	for _, no := range l.cxt.filenos {
		if err != nil {
			break
		}
		if f, e := uOpen(FT_Location, no, l.cxt); e != nil {
			return e
		} else {
			wg.Add(1)
			fq <- f
		}
	}
	wg.Wait()
	return
}

func loadFile(f *os.File, set *set.Set, cxt *context) (err error) {
	defer f.Close()
	hlen := len(locSeqHeader)
	bufferSize := 4096
	maxLocSeqSize := locationSeqSize(cxt.limitOpt.MaxKeySize)
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
				// fmt.Println(loc)
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
