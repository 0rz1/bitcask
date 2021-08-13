package bitcask

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/0rz1/bitcask/set"
)

func TestLoadFile(t *testing.T) {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal()
	}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fail()
		}
	}()
	locs := []*location{
		{fileno: 1, offset: 20, length: 12},
		{fileno: 2, offset: 20, length: 12},
		{fileno: 3, offset: 20, length: 12},
		{fileno: 4, offset: 20, length: 12},
		{fileno: 5, offset: 20, length: 12},
	}
	cxt := &context{
		path: path,
		limitOpt: LimitOption{
			MaxKeySize: 10,
		},
	}
	if f, err := uOpenAppend(FT_Location, 101, cxt); err != nil {
		t.Fatal()
	} else {
		for i, loc := range locs {
			bs := loc.makeSeqWithKey([]byte(strconv.Itoa(i)))
			if _, err := f.Write(bs); err != nil {
				t.Fatal()
			}
			bs = loc.makeSeqWithKey([]byte(strconv.Itoa(i + 10)))
			if _, err := f.Write(bs[:len(bs)-3]); err != nil {
				t.Fatal()
			}
		}
		f.Close()
	}
	f, err := uOpen(FT_Location, 101, cxt)
	if err != nil {
		t.Fatal()
	}
	set := set.New()
	if err := loadFile(f, set, cxt); err != nil {
		t.Fatal(err)
	}
	for i, loc := range locs {
		k1 := strconv.Itoa(i)
		k2 := strconv.Itoa(i + 10)
		if c, ok := set.Get(k1); !ok || loc.Compare(c) != 0 {
			t.Error()
		}
		if _, ok := set.Get(k2); ok {
			t.Error()
		}
	}
}
