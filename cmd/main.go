package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println(errors.Is(nil, errors.New("123")))
}
