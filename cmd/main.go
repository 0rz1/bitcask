package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/0rz1/bitcask"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Lack DB path")
		return
	}
	PATH := arguments[1]
	db, err := bitcask.Open(PATH)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		cmd := strings.Split(strings.TrimSpace(text), " ")
		switch cmd[0] {
		case "get":
			v, e := db.Get(cmd[1])
			if e != nil {
				fmt.Printf("Err: %v\n", e)
			} else {
				fmt.Printf("Val: %v\n", v)
			}
		case "add":
			if e := db.Add(cmd[1], cmd[2]); e != nil {
				fmt.Printf("Err: %v\n", e)
			}
		case "quit":
			return
		}
	}
}
