package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	q, err := ioutil.TempDir(".", "xxx")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(q)
		defer os.RemoveAll(q)
	}
}
