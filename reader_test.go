package bitcask

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
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
	loc := &location{
		fileno: 101,
		offset: 0,
		length: 10,
	}
	cxt := &context{path: path,
		max_filesize: 100}
	bs := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	if f, err := uOpenAppend(FT_Data, loc.fileno, cxt); err != nil {
		t.Fatal()
	} else {
		f.Write(bs)
		f.Close()
	}
	r := *newReader(cxt)
	if q, _ := r.read(loc); len(q) != loc.length || !bytes.Equal(q, bs[loc.offset:loc.offset+loc.length]) {
		t.Error(q)
	}
	loc.offset = 4
	if q, _ := r.read(loc); len(q) != loc.length || !bytes.Equal(q, bs[loc.offset:loc.offset+loc.length]) {
		t.Error(q)
	}
	loc.length = 12
	if q, _ := r.read(loc); len(q) != loc.length || !bytes.Equal(q, bs[loc.offset:loc.offset+loc.length]) {
		t.Error(q)
	}
	loc.fileno = 1
	if q, _ := r.read(loc); len(q) != 0 {
		t.Error()
	}
}
