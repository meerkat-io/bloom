package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
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
	logLevel atomic.Int32
	exit     = os.Exit // replaced in tests
)

func init() { logLevel.Store(WARN) }

// Debug logs
func Debug(v ...any) {
	if logLevel.Load() == DEBUG {
		writeLog(DEBUG, v)
	}
}

// Info logs
func Info(v ...any) {
	if logLevel.Load() <= INFO {
		writeLog(INFO, v)
	}
}

// Notice logs
func Notice(v ...any) {
	if logLevel.Load() <= NOTICE {
		writeLog(NOTICE, v)
	}
}

// Warn logs
func Warn(v ...any) {
	if logLevel.Load() <= WARN {
		writeLog(WARN, v)
	}
}

// Error logs
func Error(v ...any) {
	if logLevel.Load() <= ERROR {
		writeLog(ERROR, v)
	}
}

// Critical logs
func Critical(v ...any) {
	writeLog(CRITICAL, v)
	exit(1)
}

// Debugf logs
func Debugf(format string, v ...any) {
	if logLevel.Load() == DEBUG {
		writeFormattedLog(DEBUG, format, v)
	}
}

// Infof logs
func Infof(format string, v ...any) {
	if logLevel.Load() <= INFO {
		writeFormattedLog(INFO, format, v)
	}
}

// Noticef logs
func Noticef(format string, v ...any) {
	if logLevel.Load() <= NOTICE {
		writeFormattedLog(NOTICE, format, v)
	}
}

// Warnf logs
func Warnf(format string, v ...any) {
	if logLevel.Load() <= WARN {
		writeFormattedLog(WARN, format, v)
	}
}

// Errorf logs
func Errorf(format string, v ...any) {
	if logLevel.Load() <= ERROR {
		writeFormattedLog(ERROR, format, v)
	}
}

// Criticalf logs
func Criticalf(format string, v ...any) {
	writeFormattedLog(CRITICAL, format, v)
	exit(1)
}

// SetLevel sets level of logging
func SetLevel(level int) error {
	if level < DEBUG || level > CRITICAL {
		return fmt.Errorf("invalid log level: %d", level)
	}
	logLevel.Store(int32(level))
	return nil
}

func position() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprint(filepath.Base(file), ":", line, ":")
}

func writeLog(level int, v []any) {
	//color, now, [level], position, content " ▶ "
	log.Println(logColors[level], logStrs[level], position(), fmt.Sprint(v...), " ▶ \x1b[0m")
}

func writeFormattedLog(level int, format string, v []any) {
	//color, now, [level], position, content " ▶ "
	log.Println(logColors[level], logStrs[level], position(), fmt.Sprintf(format, v...), " ▶ \x1b[0m")

}
