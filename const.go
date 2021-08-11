package bitcask

import "errors"

type FileType int

const (
	FT_Location FileType = iota
	FT_Data
)

var ErrDuplicateOption = errors.New("duplicate type option")

var locSeqHeader = []byte{0xA0, 0xF2, 0x0B}

var defaultCacheOption = &CacheOption{Capacity: 20}
var defaultLimitOption = &LimitOption{
	MaxFileSize:  1000,
	MaxKeySize:   10,
	MaxValueSize: 100,
}
