package bitcask

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func _addFile(ft FileType, no int, cxt *context, t *testing.T) {
	path := uGetPath(ft, no, cxt)
	if f, err := os.Create(path); err != nil {
		t.Fatal()
	} else {
		f.Close()
	}
}

func _delFile(ft FileType, no int, cxt *context, t *testing.T) {
	path := uGetPath(ft, no, cxt)
	if err := os.Remove(path); err != nil {
		t.Fatal()
	}
}

func TestContextCheck(t *testing.T) {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal()
	}
	cxt := &context{path: path}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fail()
		}
	}()
	if err := cxt.check(); err != nil || cxt.maxno() != 0 {
		t.Error()
	}
	_addFile(FT_Data, 1, cxt, t)
	if err := cxt.check(); err == nil {
		t.Error()
	}
	_addFile(FT_Location, 1, cxt, t)
	if err := cxt.check(); err != nil || cxt.maxno() != 1 {
		t.Error()
	}
	_addFile(FT_Data, 2, cxt, t)
	if err := cxt.check(); err == nil {
		t.Error()
	}
	_addFile(FT_Location, 2, cxt, t)
	if err := cxt.check(); err != nil || cxt.maxno() != 2 {
		t.Error()
	}
	_addFile(FT_Data, 3, cxt, t)
	_delFile(FT_Data, 1, cxt, t)
	if err := cxt.check(); err == nil {
		t.Error()
	}
	_delFile(FT_Location, 1, cxt, t)
	_addFile(FT_Location, 3, cxt, t)
	if err := cxt.check(); err != nil || cxt.maxno() != 3 {
		t.Error()
	}

	junkpath := filepath.Join(path, "kkk")
	if f, err := os.Create(junkpath); err != nil {
		t.Fatal()
	} else {
		f.Close()
	}
	if err := cxt.check(); err == nil {
		t.Error()
	}
}
