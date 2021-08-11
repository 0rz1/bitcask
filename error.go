package bitcask

import (
	"errors"
	"fmt"
)

type Err struct {
	err error
	sub error
}

func (e *Err) Error() string {
	if e.sub == nil {
		return e.err.Error()
	}
	return fmt.Sprintf("%v: %v", e.err, e.sub)
}

func (e *Err) Unwrap() error {
	return e.err
}

var ErrDuplicateOption = errors.New("duplicate option")

var ErrKeyNotFound = errors.New("not found key")
var ErrKeyLenTooLong = errors.New("key too long")
var ErrValueLenTooLong = errors.New("value too long")

var ErrRead = errors.New("read")
var ErrWrite = errors.New("write")
var ErrNotReady = errors.New("not ready")

var ErrCxtInvalid = errors.New("cxt: other files")
var ErrCxtInconsistency = errors.New("cxt: inconsistency")
