package bitcask

import "errors"

type FileType int

const (
	FT_Invalid FileType = iota
	FT_Location
	FT_Data
)

var ErrDuplicateOption = errors.New("duplicate type option")
var ErrDiskRD = errors.New("disk read error")
var ErrDiskWR = errors.New("disk write error")
var ErrDiskUnReady = errors.New("disk unready")

var ErrCxtHasDir = errors.New("folder: has dir")
var ErrCxtInvalidName = errors.New("folder: invalid name")
var ErrCxtInconsistency = errors.New("folder: inconsistency")

var locSeqHeader = []byte{0xA0, 0xF2, 0x0B}

var defaultCacheOption = &CacheOption{Capacity: 20}
var defaultLimitOption = &LimitOption{
	MaxFileSize:  1000,
	MaxKeySize:   10,
	MaxValueSize: 100,
}
