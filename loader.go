package bitcask

import (
	"os"

	"github.com/0rz1/bitcask/set"
)

type loader struct {
	cxt *context
}

func (l *loader) load(locs *set.Set) error {
	for _, no := range l.cxt.filenos {
		f, err := uOpen(FT_Location, no, l.cxt)
		if err != nil {
			return err
		}
		if err := loadFile(f, locs); err != nil {
			return err
		}
	}
	return nil
}

func loadFile(f *os.File, locs *set.Set) error {
	return nil
}
