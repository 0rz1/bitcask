package bitcask

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAppendFile(t *testing.T) {
	path, err := ioutil.TempDir(".", "tmp")
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal()
		}
	}()
	if err != nil {
		t.Fatal()
	}
	cxt := &context{path: path, limitOpt: LimitOption{MaxFileSize: 100}}
	ap := appendFile{
		ft:  FT_Data,
		no:  10,
		cxt: cxt,
	}
	if ap.prepare(); !ap.prepared() {
		t.Fatalf("not prepared %v%v", ap.freshed, ap.aligned)
	}
	if ap.file == nil || ap.offset != 0 {
		t.Fatalf("wrong state")
	}
	if ap.exLimit(100) {
		t.Errorf("limit")
	}
	if !ap.exLimit(101) {
		t.Errorf("limit")
	}
	bs := []byte("abc")
	ap.write(bs)
	if ap.offset != 3 {
		t.Fatalf("wrong offset")
	}
	if ap.file != nil {
		ap.file.Close()
	}
}

func TestAppender(t *testing.T) {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal()
	}
	cxt := &context{path: path,
		limitOpt: LimitOption{MaxFileSize: 100}}
	app := &appender{
		cxt: cxt,
		no:  0,
		loc: &appendFile{ft: FT_Location, cxt: cxt},
		dat: &appendFile{ft: FT_Data, cxt: cxt},
	}
	defer func() {
		app.close()
		err := os.RemoveAll(path)
		if err != nil {
			t.Fail()
		}
	}()
	bskey := []byte("abc")
	bsvalue := make([]byte, 40)

	if loc, err := app.append(bskey, bsvalue); err != nil {
		t.Log(err)
		t.Fatal()
	} else if loc.fileno != 0 || loc.offset != 0 || loc.length != 40 {
		t.Fatal()
	}
	if loc, err := app.append(bskey[1:], bsvalue); err != nil {
		t.Fatal()
	} else if loc.fileno != 0 || loc.offset != 40 || loc.length != 40 {
		t.Fatal()
	}
	if loc, err := app.append(bskey[2:], bsvalue[:30]); err != nil {
		t.Fatal()
	} else if loc.fileno != 1 || loc.offset != 0 || loc.length != 30 {
		t.Fatal()
	}
}
