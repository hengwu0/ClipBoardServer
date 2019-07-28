//总日志
package algorithm

import (
	"log"
	"os"
	"time"
)

var file *os.File
var logger *log.Logger
var logdir = "./log/"

func init() {
	var err error
	if err = os.MkdirAll(logdir, 0755); err != nil {
		log.Fatalln(`fail to create dir: "./log"! %v`, err)
	}
	filename := logdir + "ClipBoard_" + time.Now().Format("2006-01-02") + ".log"
	if file, err = os.Create(filename); err != nil {
		log.Fatalln("fail to create %s file! %v", filename, err)
	}
	logger = log.New(file, "", log.LstdFlags) //log.Lshortfile
}

func Close() {
	file.Close()
}

func (alist *Alist) GetLogger() *log.Logger {
	if alist.logger == nil {
		filename := logdir + alist.hash + "_" + time.Now().Format("2006-01-02") + ".log"
		if f, err := os.Create(filename); err == nil {
			alist.file = f
			alist.logger = log.New(alist.file, "", log.LstdFlags)
			return alist.logger
		} else {
			logger.Printf("fail to create：(%s)%v %v", filename, err, f)
		}
		alist.logger = logger
		return logger
	}
	return alist.logger
}
