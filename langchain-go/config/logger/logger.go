package logger

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	INFO
	WARNING
	ERROR
	DISABLED
)

var (
	logger    *log.Logger
	logOutput io.Writer = os.Stderr
	logLevel  LogLevel  = INFO
)

func init() {
	updateLogger()
}

func updateLogger() {
	logger = log.New(logOutput, "", log.LstdFlags)
}

func SetOutput(w io.Writer) {
	logOutput = w
	updateLogger()
}

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func Trace(v ...interface{}) {
	if logLevel <= TRACE {
		logger.SetPrefix("[TRACE] ")
		logger.Println(v...)
	}
}

func Info(v ...interface{}) {
	if logLevel <= INFO {
		logger.SetPrefix("[INFO] ")
		logger.Println(v...)
	}
}

func Warning(v ...interface{}) {
	if logLevel <= WARNING {
		logger.SetPrefix("[WARNING] ")
		logger.Println(v...)
	}
}

func Error(v ...interface{}) {
	if logLevel <= ERROR {
		logger.SetPrefix("[ERROR] ")
		logger.Println(v...)
	}
}

func Disabled(v ...interface{}) {
	if logLevel <= DISABLED {
		logger.SetPrefix("[DISABLED] ")
		logger.Println(v...)
	}
}
