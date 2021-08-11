package bitcask

type FileType int

const (
	FT_Invalid FileType = iota
	FT_Location
	FT_Data
)

var locSeqHeader = []byte{0xA0, 0xF2, 0x0B}

var defaultCacheOption = &CacheOption{Capacity: 20}
var defaultLimitOption = &LimitOption{
	MaxFileSize:  1000,
	MaxKeySize:   10,
	MaxValueSize: 100,
}
var defaultDiskOption = &DiskOption{
	ReaderCnt: 4,
	LoaderCnt: 4,
}
