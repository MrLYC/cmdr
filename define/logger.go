package define

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cast"
	adapter "logur.dev/adapter/template"
	"logur.dev/logur"
)

var Logger logur.Logger

type terminalLogger struct {
	adapter.Logger
	level          logur.Level
	withErrorStack bool
	traceStyle     color.Style
	debugStyle     color.Style
	infoStyle      color.Style
	warnStyle      color.Style
	errorStyle     color.Style
}

func (l *terminalLogger) getMessages(msg string) string {
	return strings.ToUpper(msg[:1]) + msg[1:]
}

func (l *terminalLogger) getMessagesByFields(fields []map[string]interface{}) []string {
	var (
		messages   = []string{}
		errorKey   string
		errorValue error
	)

	if len(fields) == 0 {
		return messages
	}

	for k, v := range fields[0] {
		switch value := v.(type) {
		case error:
			errorKey = k
			errorValue = value
		default:
			messages = append(messages, fmt.Sprintf("%s=%s", k, cast.ToString(v)))
		}
	}

	if errorValue == nil {
		return messages
	}

	if l.withErrorStack {
		messages = append(messages, fmt.Sprintf("%s=%+v", errorKey, errorValue))
	} else {
		messages = append(messages, fmt.Sprintf("%s=%v", errorKey, errorValue))
	}

	return messages
}

func (l *terminalLogger) log(level logur.Level, style color.Style, msg string, fields []map[string]interface{}) {
	if level < l.level || msg == "" {
		return
	}

	messages := []string{l.getMessages(msg)}
	messages = append(messages, l.getMessagesByFields(fields)...)

	style.Println(strings.Join(messages, ", "))
}

// Trace implements the Logur Logger interface.
func (l *terminalLogger) Trace(msg string, fields ...map[string]interface{}) {
	l.log(logur.Trace, l.traceStyle, msg, fields)
}

// Debug implements the Logur Logger interface.
func (l *terminalLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(logur.Debug, l.debugStyle, msg, fields)
}

// Info implements the Logur Logger interface.
func (l *terminalLogger) Info(msg string, fields ...map[string]interface{}) {
	l.log(logur.Info, l.infoStyle, msg, fields)
}

// Warn implements the Logur Logger interface.
func (l *terminalLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(logur.Warn, l.warnStyle, msg, fields)
}

// Error implements the Logur Logger interface.
func (l *terminalLogger) Error(msg string, fields ...map[string]interface{}) {
	l.log(logur.Error, l.errorStyle, msg, fields)
}

func InitLogger() {
	level, ok := logur.ParseLevel(Configuration.GetString("log.level"))
	if !ok {
		level = logur.Info
	}

	Logger = &terminalLogger{
		level:          level,
		withErrorStack: level < logur.Info,
		traceStyle:     color.Style{color.OpItalic, color.FgDarkGray},
		debugStyle:     color.Style{color.OpItalic, color.FgDefault},
		infoStyle:      color.Style{color.OpBold, color.FgDefault},
		warnStyle:      color.Style{color.OpBold, color.FgYellow},
		errorStyle:     color.Style{color.OpBold, color.FgRed},
	}
}

func init() {
	Logger = logur.NewNoopLogger()
}
