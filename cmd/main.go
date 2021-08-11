package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"

	"github.com/0rz1/bitcask"
)

func makerand(len int) []byte {
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		bs[i] = byte(rand.Intn(127))
	}
	return bs
}

func foo() {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		fmt.Println(err)
		return
	}
	// db, err := bitcask.Open(path, &bitcask.DiskOption{ReaderCnt: 3, LoaderCnt: 1})
	db, err := bitcask.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	length := 60
	bs := makerand(length)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; ; j++ {
				r := rand.Intn(10)
				k := strconv.Itoa(rand.Intn(1000))
				if r == 0 {
					_, err := db.Get(k)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					p := rand.Intn(length - 10)
					v := bs[p:]
					err := db.Add(k, string(v))
					if err != nil {
						fmt.Println(err)
					}
				}
				if j%10 == 0 {
					fmt.Printf("p %v: %v\n", i, j)
				}
				if j == 100 {
					break
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func foo1() {
	path, err := ioutil.TempDir(".", "tmp")
	defer os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = bitcask.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// length := 60
}

func main() {
	foo()

}
