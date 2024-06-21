package internal

import (
	"fmt"
	"io"
	"os"
)

type (
	LogLevel     int
	LogFormatter func(severity LogLevel, v ...any) string
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

type RummageLogger interface {
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Err(v ...any)
	Fatal(v ...any)
}

type GeneralRummageLoggerImpl struct {
	Formatter LogFormatter // formatter func
	Output    io.Writer    // write output to the io.writer, usually os.stdout
}

func NewLogger(formatter LogFormatter, output io.Writer) *GeneralRummageLoggerImpl {
	if formatter == nil {
		formatter = defaultFormatter
	}
	return &GeneralRummageLoggerImpl{
		Formatter: formatter,
		Output:    output,
	}
}

func (l *GeneralRummageLoggerImpl) Debug(v ...any) {
	l.log(DEBUG, v...)
}

func (l *GeneralRummageLoggerImpl) Info(v ...any) {
	l.log(INFO, v...)
}

func (l *GeneralRummageLoggerImpl) Warn(v ...any) {
	l.log(WARN, v...)
}

func (l *GeneralRummageLoggerImpl) Err(v ...any) {
	l.log(ERROR, v...)
}

func (l *GeneralRummageLoggerImpl) Fatal(v ...any) {
	l.log(FATAL, v...)
	os.Exit(1)
}

// General log that writes input to the io.writer
func (l *GeneralRummageLoggerImpl) log(severity LogLevel, v ...any) {
	b := []byte(l.Formatter(severity, v...))
	_, err := l.Output.Write(b)
	if err != nil {
		panic(err)
	}
}

// RummageLogger's default formatter
func defaultFormatter(severity LogLevel, v ...any) string {
	level := [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[severity]
	color := LogLevelColors[severity]
	msg := fmt.Sprint(v...)
	return fmt.Sprintf("%s%s%s | %s%s", color, level, RESET, msg, "\n")
}
