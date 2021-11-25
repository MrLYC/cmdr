package core

import (
	"context"

	"github.com/mrlyc/cmdr/define"
)

type StepLogger struct {
	BaseStep
	message     string
	contextKeys []define.ContextKey
}

func (s *StepLogger) String() string {
	return "step-logger"
}

func (s *StepLogger) Run(ctx context.Context) (context.Context, error) {
	fields := make(map[string]interface{}, len(s.contextKeys))

	for _, key := range s.contextKeys {
		fields[key.String()] = ctx.Value(key)
	}

	define.Logger.Info(s.message, fields)
	return ctx, nil
}

func NewStepLoggerWithFields(message string, keys ...define.ContextKey) *StepLogger {
	return &StepLogger{
		message:     message,
		contextKeys: keys,
	}
}

func NewStepLogger(message string) *StepLogger {
	return NewStepLoggerWithFields(message)
}
