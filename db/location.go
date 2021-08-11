package db

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sync"
	"sync/atomic"
	"unsafe"
)

var locationHeader = []byte{0xA0, 0xF2, 0x0B}

type location struct {
	fileno int
	offset int
	length int
}

func makeLocationAndKey(bs []byte) (loc location, key []byte, ok bool) {
	//header + location(3*4) + len(4) + key + crc(4)
	headerlen := len(locationHeader)
	nokeySize := locationSeqSize(0)
	if len(bs) <= nokeySize || !bytes.Equal(bs[:headerlen], locationHeader) {
		return
	}
	keylen := len(bs) - nokeySize
	if int(binary.BigEndian.Uint32(bs[headerlen+12:headerlen+16])) != keylen {
		return
	}
	checksum := binary.BigEndian.Uint32(bs[len(bs)-4:])
	if checksum != crc32.ChecksumIEEE(bs[:len(bs)-4]) {
		return
	}
	ok = true
	loc.fileno = int(binary.BigEndian.Uint32(bs[headerlen : headerlen+4]))
	loc.offset = int(binary.BigEndian.Uint32(bs[headerlen+4 : headerlen+8]))
	loc.length = int(binary.BigEndian.Uint32(bs[headerlen+8 : headerlen+12]))
	key = make([]byte, keylen)
	copy(key, bs[headerlen+16:])
	return
}

func locationSeqSize(ksz int) int {
	//header + location(3*4) + len(4) + key + crc(4)
	return len(locationHeader) + 12 + 4 + ksz + 4
}

func (l *location) MakeSeqWithKey(key []byte) []byte {
	//header + location(3*4) + len(4) + key + crc(4)
	size := locationSeqSize(len(key))
	bs := make([]byte, size)
	copy(bs, locationHeader)
	off := len(locationHeader)
	binary.BigEndian.PutUint32(bs[off:], uint32(l.fileno))
	binary.BigEndian.PutUint32(bs[off+4:], uint32(l.offset))
	binary.BigEndian.PutUint32(bs[off+8:], uint32(l.length))
	binary.BigEndian.PutUint32(bs[off+12:], uint32(len(key)))
	copy(bs[off+16:], key)
	binary.BigEndian.PutUint32(bs[size-4:], crc32.ChecksumIEEE(bs[:size-4]))
	return bs
}

func (l *location) Compare(a *location) int {
	if l.fileno != a.fileno {
		return l.fileno - a.fileno
	}
	return l.offset - a.offset
}

type locationMap struct {
	store *sync.Map
}

func newLocationMap() *locationMap {
	return &locationMap{
		store: &sync.Map{},
	}
}

func (mp *locationMap) add(key string, loc *location) {
	locptr := unsafe.Pointer(loc)
	locwrapper := &struct{ v unsafe.Pointer }{v: locptr}
	act, loaded := mp.store.LoadOrStore(key, locwrapper)
	if loaded {
		actw := act.(*struct{ v unsafe.Pointer })
		for {
			pv := actw.v
			if (*location)(pv).Compare(loc) <= 0 {
				return
			}
			if atomic.CompareAndSwapPointer(&actw.v, pv, locptr) {
				return
			}
		}
	}
}

func (mp *locationMap) get(key string) (*location, bool) {
	val, ok := mp.store.Load(key)
	if !ok {
		return nil, false
	}
	valw := val.(*struct{ v unsafe.Pointer })
	vptr := atomic.LoadPointer(&valw.v)
	return (*location)(vptr), ok
}
