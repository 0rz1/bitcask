package bitcask

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"

	"github.com/0rz1/bitcask/set"
)

type location struct {
	fileno int
	offset int
	length int
}

var _ set.Comparable = &location{}

func locationSeqSize(ksz int) int {
	//header + location(3*4) + len(4) + key + crc(4)
	return len(locSeqHeader) + 12 + 4 + ksz + 4
}

func makeLocationAndKey(bs []byte) (loc *location, key []byte, ok bool) {
	loc = &location{}
	//header + location(3*4) + len(4) + key + crc(4)
	headerlen := len(locSeqHeader)
	nokeySize := locationSeqSize(0)
	if len(bs) <= nokeySize || !bytes.Equal(bs[:headerlen], locSeqHeader) {
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

func (l *location) makeSeqWithKey(key []byte) []byte {
	//header + location(3*4) + len(4) + key + crc(4)
	size := locationSeqSize(len(key))
	bs := make([]byte, size)
	copy(bs, locSeqHeader)
	off := len(locSeqHeader)
	binary.BigEndian.PutUint32(bs[off:], uint32(l.fileno))
	binary.BigEndian.PutUint32(bs[off+4:], uint32(l.offset))
	binary.BigEndian.PutUint32(bs[off+8:], uint32(l.length))
	binary.BigEndian.PutUint32(bs[off+12:], uint32(len(key)))
	copy(bs[off+16:], key)
	binary.BigEndian.PutUint32(bs[size-4:], crc32.ChecksumIEEE(bs[:size-4]))
	return bs
}

func (loc *location) Compare(other set.Comparable) int {
	if o, ok := other.(*location); ok {
		if loc.fileno == o.fileno {
			return loc.offset - o.offset
		}
		return loc.fileno - o.fileno
	}
	return 0
}
