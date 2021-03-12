package logger

import (
	//"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// Level type
type Level int

// Log level, from low to high, more high means more serious
const (
	LevelFatal Level = iota
	LevelError
	LevelInfo
	LevelDebug
)

// LogLevel sets log level
//var LogLevel = LevelError

// LogPath set log file path
//var LogPath string
var logger = NewLogger()

// Logger represents logging object
type Logger struct {
	level           Level
	logpath         string
	errorstacktrace bool
}

// NewLogger creates default logger
func NewLogger() *Logger {
	var l = new(Logger)
	l.level = LevelError

	return l
}

// SetLevel changes the logger level
func SetLevel(level Level) {
	//LogLevel = level
	logger.SetLevel(level)
}

// SetLevel sets log level, any log level less than it will not log
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetLogPath sets log file name
func SetLogPath(logfile string) {
	//LogPath = logfile
	logger.SetLogPath(logfile)
}

// SetLogPath sets log file name
func (l *Logger) SetLogPath(logfile string) {
	l.logpath = logfile
}

// EnableErrorStackTrace enables stack trace for ErrorTracef method
func EnableErrorStackTrace() {
	logger.EnableErrorStackTrace()
}

// EnableErrorStackTrace enables stack trace for ErrorTracef method
func (l *Logger) EnableErrorStackTrace() {
	l.errorstacktrace = true
}

// Fatalf records the log with fatal level and exits
func Fatalf(format string, args ...interface{}) {
	logger.Output(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Errorf records the log with error level
func Errorf(format string, args ...interface{}) {
	logger.Output(LevelError, fmt.Sprintf(format, args...))
}

// Infof records the log with info level
func Infof(format string, args ...interface{}) {
	logger.Output(LevelInfo, fmt.Sprintf(format, args...))
}

// Debugf records the log with debug level
func Debugf(format string, args ...interface{}) {
	logger.Output(LevelDebug, fmt.Sprintf(format, args...))
}

// Errorf records the log with error level
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, args...))
}

// Fatalf records the log with fatal level and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Infof records the log with info level
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, args...))
}

// Debugf records the log with debug level
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, args...))
}

// ErrorTracef records the log with stack trace in error level
func ErrorTracef(format string, args ...interface{}) {
	logger.ErrorTracef(format, args...)
}

// ErrorTracef records the log with stack trace in error level
func (l *Logger) ErrorTracef(format string, args ...interface{}) {
	if l.errorstacktrace {
		strace := errors.New(fmt.Sprintf(format, args...))
		l.Output(LevelError, fmt.Sprintf("%+v", strace))
	} else {
		l.Output(LevelError, fmt.Sprintf(format, args...))
	}
}

// Output records the log with special callstack depth and log level.
func (l *Logger) Output(level Level, msg string) {
	//if l.level < level && LogLevel < level {
	if l.level < level {
		return
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	if l.logpath != "" {
		logf, err := os.OpenFile(l.logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer logf.Close()
		log.SetOutput(logf)
	}

	var logLevelStr = "[UNKNOWN]"
	switch level {
	case LevelFatal:
		logLevelStr = "[FATAL]"
	case LevelError:
		logLevelStr = "[ERROR]"
	case LevelInfo:
		logLevelStr = "[INFO ]"
	case LevelDebug:
		logLevelStr = "[DEBUG]"
	}

	log.Printf("%s %s:%d %s: %s", logLevelStr, filepath.Base(file), line, fnName, msg)
}
