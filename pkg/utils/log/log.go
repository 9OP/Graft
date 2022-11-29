package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"graft/pkg/utils"
)

func ConfigureLogger(id string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix)
	logFile, err := os.OpenFile(fmt.Sprintf("%s/log_%s.txt", utils.GraftPath(), id), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

type LogLevel uint8

const (
	_ LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var level LogLevel = DEBUG

func Debugf(format string, args ...interface{}) {
	if level > DEBUG {
		return
	}
	log.SetPrefix(fmt.Sprintf("%-10s", "[DEBUG]"))
	log.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	if level > INFO {
		return
	}
	log.SetPrefix(fmt.Sprintf("%-10s", "[INFO]"))
	log.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	if level > WARN {
		return
	}
	log.SetPrefix(fmt.Sprintf("%-10s", "[WARN]"))
	log.Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if level > ERROR {
		return
	}
	log.SetPrefix(fmt.Sprintf("%-10s", "[ERROR]"))
	log.Printf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.SetPrefix(fmt.Sprintf("%-10s", "[FATAL]"))
	log.Fatalf(format, args...)
}
