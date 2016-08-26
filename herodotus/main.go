package herodotus

import (
	"log"
	"os"
)

func CreateFileLog(filename string) *log.Logger {
	if f, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		panic(err)
	} else {
		// file never gets closed in the program, hopefully the OS pulls its weight
		return log.New(f, "", log.LstdFlags)
	}
}
