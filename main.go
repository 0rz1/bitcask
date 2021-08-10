package main

import (
	"fmt"
	"math/rand"
)

func main() {
	for i := 0; i < 100; i++ {
		a, b := rand.Int(), rand.Int()
		c, d := a, b
		a &^= b
		c &= ^d
		if c != a {
			fmt.Println("bad")
		}
	}
	fmt.Println("done")
}
