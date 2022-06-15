package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// log levels
const (
	_ = iota
	DEBUG
	INFO
	NOTICE
	WARN
	ERROR
	CRITICAL
)

var (
	logColors = []string{
		"",
		"\x1b[36;1m",
		"\x1b[0m",
		"\x1b[32m",
		"\x1b[33m",
		"\x1b[31m",
		"\x1b[31;1m",
	}
	logStrs = []string{
		"",
		"[DEBUG]   ",
		"[INFO]    ",
		"[NOTICE]  ",
		"[WARN]    ",
		"[ERROR]   ",
		"[CRITICAL]",
	}
	logLevel = WARN
)

// Debug logs
func Debug(v ...interface{}) {
	if logLevel == DEBUG {
		writeLog(DEBUG, v)
	}
}

// Info logs
func Info(v ...interface{}) {
	if logLevel <= INFO {
		writeLog(INFO, v)
	}
}

// Notice logs
func Notice(v ...interface{}) {
	if logLevel <= NOTICE {
		writeLog(NOTICE, v)
	}
}

// Warn logs
func Warn(v ...interface{}) {
	if logLevel <= WARN {
		writeLog(WARN, v)
	}
}

// Error logs
func Error(v ...interface{}) {
	if logLevel <= ERROR {
		writeLog(ERROR, v)
	}
}

// Critical logs
func Critical(v ...interface{}) {
	writeLog(CRITICAL, v)
	os.Exit(1)
}

// Debugf logs
func Debugf(format string, v ...interface{}) {
	if logLevel == DEBUG {
		writeFormattedLog(DEBUG, format, v)
	}
}

// Infof logs
func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		writeFormattedLog(INFO, format, v)
	}
}

// Noticef logs
func Noticef(format string, v ...interface{}) {
	if logLevel <= NOTICE {
		writeFormattedLog(NOTICE, format, v)
	}
}

// Warnf logs
func Warnf(format string, v ...interface{}) {
	if logLevel <= WARN {
		writeFormattedLog(WARN, format, v)
	}
}

// Errorf logs
func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		writeFormattedLog(ERROR, format, v)
	}
}

// Criticalf logs
func Criticalf(format string, v ...interface{}) {
	writeFormattedLog(CRITICAL, format, v)
	os.Exit(1)
}

// SetLevel sets level of logging
func SetLevel(level int) error {
	if level < DEBUG || level > CRITICAL {
		return fmt.Errorf("invalid log level: %d", level)
	}
	logLevel = level
	return nil
}

func position() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprint(filepath.Base(file), ":", line, ":")
}

func writeLog(level int, v []interface{}) {
	//color, now, [level], position, content " ▶ "
	log.Println(logColors[level], logStrs[level], position(), fmt.Sprint(v...), " ▶ \x1b[0m")
}

func writeFormattedLog(level int, format string, v []interface{}) {
	//color, now, [level], position, content " ▶ "
	log.Println(logColors[level], logStrs[level], position(), fmt.Sprintf(format, v...), " ▶ \x1b[0m")

}
