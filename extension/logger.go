package extension

import (
	"fmt"
	"io"
)

const NameForLog string = "Invocation History Extension"

type Logger struct {
	prefix string
	writer io.Writer
}

func NewLogger(w io.Writer, prefix string) *Logger {
	p := ""
	if len(prefix) > 0 {
		p = fmt.Sprintf("[%s] ", prefix)
	}

	return &Logger{
		prefix: p,
		writer: w,
	}
}

func (l *Logger) Info(format string, v ...any) {
	l.write(logLevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...any) {
	l.write(logLevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...any) {
	l.write(logLevelError, fmt.Sprintf(format, v...))
}

type logLevel string

const (
	logLevelInfo  logLevel = "INFO"
	logLevelWarn  logLevel = "WARN"
	logLevelError logLevel = "ERROR"
)

func (lv logLevel) string() string {
	return string(lv)
}

func (l *Logger) write(level logLevel, m string) {

	fmt.Fprintf(l.writer, "%s%s: %s\n", l.prefix, level.string(), m)
}
