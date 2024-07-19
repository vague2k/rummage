package logger

import (
	"fmt"
	"io"
	"os"
)

const (
	RESET  = "\033[1;0m"
	RED    = "\033[1;31m"
	GREEN  = "\033[1;32m"
	YELLOW = "\033[1;33m"
	PURPLE = "\033[1;35m"
	CYAN   = "\033[1;36m"
	WHITE  = "\033[1;37m"
)

type (
	LogLevel int
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var LogLevelColors = map[LogLevel]string{
	DEBUG: CYAN,
	INFO:  GREEN,
	WARN:  YELLOW,
	ERROR: RED,
	FATAL: PURPLE,
}

type Logger struct {
	Output io.Writer // write output to the io.writer, usually os.stdout
}

func New() *Logger {
	return &Logger{
		Output: os.Stdout,
	}
}

func (l *Logger) Debug(v ...any) {
	l.log(DEBUG, v...)
}

func (l *Logger) Info(v ...any) {
	l.log(INFO, v...)
}

func (l *Logger) Warn(v ...any) {
	l.log(WARN, v...)
}

func (l *Logger) Err(v ...any) {
	l.log(ERROR, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.log(FATAL, v...)
	os.Exit(1)
}

// General log that writes input to the io.writer
func (l *Logger) log(severity LogLevel, v ...any) {
	b := []byte(formatLog(severity, v...))
	_, err := l.Output.Write(b)
	if err != nil {
		panic(err)
	}
}

func formatLog(severity LogLevel, v ...any) string {
	level := [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[severity]
	color := LogLevelColors[severity]
	msg := fmt.Sprint(v...)
	return fmt.Sprintf("%s%s%s | %s%s", color, level, RESET, msg, "\n")
}
