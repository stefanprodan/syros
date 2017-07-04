package main

import (
	"io"
	"log"
	"os"
	"path"
)

func setLogFile(dir string) {
	logpath := path.Join(dir, "deployctl.log")
	logFile, err := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	logFile.Truncate(0)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(0)
	log.SetOutput(mw)
}
