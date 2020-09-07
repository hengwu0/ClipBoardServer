package main

import (
	"Algorithm"
	"Protocol"
	"flag"
	"fmt"
	"os"
)

func writePid() {
	pid := os.Getpid()
	if f, err := os.Create("pidfile"); err != nil {
		fmt.Fprintf(os.Stderr, "Can't write pidfile: %v\n", err)
		os.Exit(1)
	} else {
		f.WriteString(fmt.Sprint(pid))
		f.Close()
	}

}

func main() {
	var h bool
	var c string
	flag.BoolVar(&h, "h", false, "show this help")
	flag.StringVar(&c, "c", "", "set `CONFIG_PATH` path")
	flag.Parse()

	if c != "" {
		os.Setenv("CONFIG_PATH", c)
	}
	if h || flag.NArg() > 0 {
		flag.PrintDefaults()
		return
	}
	writePid()
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
