package bitcask

import "errors"

var ErrDuplicateOption = errors.New("duplicate type option")

var defaultCacheOption = &CacheOption{Capacity: 20}
var defaultLimitOption = &LimitOption{
	MaxFileSize:  1000,
	MaxKeySize:   10,
	MaxValueSize: 100,
}
