package bitcask

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func makerandbs(len int) []byte {
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		bs[i] = byte(rand.Intn(127))
	}
	return bs
}

func randupdatedb(db *DB, mp map[string]string, sz int, rlimit int) error {
	length := 60
	bs := makerandbs(length)
	for j := 0; j < sz*2; j++ {
		k := strconv.Itoa(rand.Intn(sz))
		r := rand.Intn(rlimit)
		if r != 0 {
			_, err := db.Get(k)
			if err != nil && err != ErrKeyNotFound {
				return err
			}
		} else {
			p := rand.Intn(length - 10)
			v := bs[p:]
			err := db.Add(k, string(v))
			if err != nil {
				return err
			}
			if mp != nil {
				mp[k] = string(v)
			}
		}
	}
	return nil
}

func TestDB(t *testing.T) {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal(err)
		return
	}
	db, err := Open(path)
	defer func() {
		if db != nil {
			db.Close()
			os.RemoveAll(path)
		}
	}()
	if err != nil {
		t.Fatal(err)
		return
	}
	mp := make(map[string]string)
	if err := randupdatedb(db, mp, 100, 1); err != nil {
		t.Fatal(err)
	}
	for k, v := range mp {
		if rv, err := db.Get(k); err != nil || rv != v {
			t.Error()
		}
	}
	db.Close()
	db, err = Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range mp {
		if rv, err := db.Get(k); err != nil || rv != v {
			t.Error()
		}
	}
	db.Close()
}

func BenchmarkDB(b *testing.B) {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		b.Fatal(err)
		return
	}
	db, err := Open(path)
	defer func() {
		if db != nil {
			db.Close()
			os.RemoveAll(path)
		}
	}()
	if err != nil {
		b.Fatal(err)
		return
	}
	if err := randupdatedb(db, nil, 200, 1); err != nil {
		b.Fatal(err)
	}
	length := 100
	bs := makerandbs(length)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			r := rand.Intn(2)
			k := strconv.Itoa(rand.Intn(200))
			if r == 0 {
				p := rand.Intn(length)
				v := bs[p:]
				db.Add(k, string(v))
			} else {
				db.Get(k)
			}
		}
	})
}
