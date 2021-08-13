package bitcask

import (
	"io/ioutil"
	"sort"
)

type context struct {
	path     string
	limitOpt LimitOption
	diskOpt  DiskOption
	filenos  []int
}

func (c *context) maxno() int {
	if len(c.filenos) == 0 {
		return 0
	}
	return c.filenos[len(c.filenos)-1]
}

func (c *context) check() error {
	fs, err := ioutil.ReadDir(c.path)
	if err != nil {
		return err
	}
	dats := []int{}
	locs := []int{}
	for _, f := range fs {
		if f.IsDir() {
			return ErrCxtInvalid
		}
		ft, no := uGetFTAndNo(f.Name())
		if ft == FT_Invalid || no < 0 {
			return ErrCxtInvalid
		} else if ft == FT_Data {
			dats = append(dats, no)
		} else if ft == FT_Location {
			locs = append(locs, no)
		} else {
			panic("unknown error")
		}
	}
	if len(dats) != len(locs) {
		return ErrCxtInconsistency
	}
	sort.Ints(dats)
	sort.Ints(locs)
	for i := range dats {
		if dats[i] != locs[i] {
			return ErrCxtInconsistency
		}
	}
	c.filenos = dats
	return nil
}
