package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
	"github.com/spf13/cast"
	adapter "logur.dev/adapter/template"
	"logur.dev/logur"
)

var globalLogger logur.Logger

type terminalLogger struct {
	adapter.Logger
	level          logur.Level
	withErrorStack bool
	errorKey       string
}

func (l *terminalLogger) getFieldsMessages(fields []map[string]interface{}) []string {
	var (
		messages   = []string{}
		errorValue error
	)

	if len(fields) == 0 {
		return messages
	}

	for k, v := range fields[0] {
		if k != l.errorKey {
			messages = append(messages, fmt.Sprintf("%s=%s", k, cast.ToString(v)))
		} else if v != nil {
			errorValue = v.(error)
		}
	}

	if errorValue != nil {
		if l.withErrorStack {
			messages = append(messages, fmt.Sprintf("%s=%+v", l.errorKey, errorValue))
		} else {
			messages = append(messages, fmt.Sprintf("%s=%v", l.errorKey, errorValue))
		}
	}

	return messages
}

func (l *terminalLogger) log(level logur.Level, msg string, fields []map[string]interface{}, fn func(message string) fmt.Stringer) {
	if level < l.level || msg == "" {
		return
	}

	messages := []string{strings.ToUpper(msg[:1]) + msg[1:]}
	messages = append(messages, l.getFieldsMessages(fields)...)

	fmt.Fprintln(os.Stderr, fn(strings.Join(messages, ", ")))
}

// Trace implements the Logur Logger interface.
func (l *terminalLogger) Trace(msg string, fields ...map[string]interface{}) {
	l.log(logur.Trace, msg, fields, func(message string) fmt.Stringer {
		return termenv.String(message).Italic().Underline()
	})
}

// Debug implements the Logur Logger interface.
func (l *terminalLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(logur.Debug, msg, fields, func(message string) fmt.Stringer {
		return termenv.String(message).Italic()
	})
}

// Info implements the Logur Logger interface.
func (l *terminalLogger) Info(msg string, fields ...map[string]interface{}) {
	l.log(logur.Info, msg, fields, func(message string) fmt.Stringer {
		return termenv.String(message).Bold()
	})
}

// Warn implements the Logur Logger interface.
func (l *terminalLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(logur.Warn, msg, fields, func(message string) fmt.Stringer {
		return termenv.String(message).Bold().Foreground(termenv.ANSIBrightYellow)
	})
}

// Error implements the Logur Logger interface.
func (l *terminalLogger) Error(msg string, fields ...map[string]interface{}) {
	l.log(logur.Error, msg, fields, func(message string) fmt.Stringer {
		return termenv.String(message).Bold().Foreground(termenv.ANSIBrightRed)
	})
}

func InitTerminalLogger(level logur.Level, withErrorStack bool, errorKey string) {
	globalLogger = &terminalLogger{
		level:          level,
		withErrorStack: withErrorStack,
		errorKey:       errorKey,
	}
}

func SetLogger(logger logur.Logger) {
	globalLogger = logger
}

func GetLogger() logur.Logger {
	return globalLogger
}

func init() {
	globalLogger = logur.NoopLogger{}
}
