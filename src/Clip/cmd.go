package main

import (
	"Algorithm"
	"Protocol"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Println("don't need more args:", os.Args[1:])
		return
	}
	if err := protocol.Listen(); err != nil {
		fmt.Fprintf(os.Stderr, "listen ERROR: %v\n", err)
		return
	}
	for {
		if conn, err := protocol.Accept(); err == nil {
			go func() {
				if node := algorithm.NewClient(conn); node != nil {
					node.ParseCmd()
				}
			}()
		} else {
			fmt.Fprintf(os.Stderr, "accept error: %s", err)
		}
	}
}
