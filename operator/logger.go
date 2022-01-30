package operator

import (
	"context"

	"github.com/mrlyc/cmdr/define"
)

type OperatorLogger struct {
	BaseOperator
	message     string
	contextKeys []define.ContextKey
}

func (s *OperatorLogger) String() string {
	return "operator-logger"
}

func (s *OperatorLogger) Run(ctx context.Context) (context.Context, error) {
	fields := make(map[string]interface{}, len(s.contextKeys))

	for _, key := range s.contextKeys {
		fields[key.String()] = ctx.Value(key)
	}

	define.Logger.Info(s.message, fields)
	return ctx, nil
}

func NewOperatorLoggerWithFields(message string, keys ...define.ContextKey) *OperatorLogger {
	return &OperatorLogger{
		message:     message,
		contextKeys: keys,
	}
}

func NewOperatorLogger(message string) *OperatorLogger {
	return NewOperatorLoggerWithFields(message)
}
