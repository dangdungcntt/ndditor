// Package logger provides a logger for the editor
package logger

import (
	"fmt"
	"log"
	"os"
)

var logFile *os.File

// WriteLog writes a message to the log file
func WriteLog(args ...any) {
	if logFile == nil {
		var err error
		logFile, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, _ = fmt.Fprintln(logFile, args...)
}
