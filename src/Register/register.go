package register

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

//ok=1:Need Update
//ok=2:Key err
func CheckVandK(logger *log.Logger, v int, k string) (int, string) {
	if v < 3157553 { //"1.0"
		return 1, ""
	}
	var index int
	if index = strings.Index(k, "--"); index < 0 || index+2 >= len(k) {
		return 2, ""
	}
	hash := k[index+2 : len(k)-1]
	if !IsRegister(logger, []byte(hash)) {
		return 2, ""
	}
	return 0, hash
}

var path = "/home/wuheng/test/ClipBoard/"
var name = "register.dat"

func IsRegister(logger *log.Logger, key []byte) bool {

	var f *os.File
	var err error
	if f, err = os.Open(path + name); err != nil {
		logger.Printf("Can't open register file: %v\n", err)
		return false
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line []byte
	comm := []byte("# ")

	for {
		if line, _, err = r.ReadLine(); err != nil {
			break
		}
		if bytes.HasPrefix(line, comm) {
			continue
		}
		if bytes.Equal(line, key) {
			logger.Printf("Register key(%s) ok!\n", key)
			return true
		}
	}
	logger.Printf("Can't find key %s(%v) in register file!\n", key, key)
	return false
}
