package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/0rz1/bitcask"
)

func main() {
	arguments := os.Args
	if len(arguments) <= 2 {
		fmt.Println("Please provide port number & DB path")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	PATH := arguments[2]
	db, err := bitcask.Open(PATH)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	for {
		reader := bufio.NewReader(c)
		text, _ := reader.ReadString('\n')
		cmd := strings.Split(strings.TrimSpace(text), " ")
		switch cmd[0] {
		case "get":
			v, e := db.Get(cmd[1])
			if e != nil {

				c.Write([]byte(fmt.Sprintf("Err: %v\n", e)))
			} else {
				c.Write([]byte(fmt.Sprintf("Val: %v\n", v)))
			}
		case "add":
			if e := db.Add(cmd[1], cmd[2]); e != nil {
				c.Write([]byte(fmt.Sprintf("Err: %v\n", e)))
			}
		case "quit":
			c.Write([]byte("quit\n"))
			return
		default:
			c.Write([]byte("\n"))
		}
	}
}
