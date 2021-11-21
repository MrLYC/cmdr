package utils

import (
	"context"

	"github.com/mrlyc/cmdr/define"
)

func GetStringFromContext(ctx context.Context, key define.ContextKey) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

func GetBoolFromContext(ctx context.Context, key define.ContextKey) bool {
	value := ctx.Value(key)
	if value == nil {
		return false
	}
	return value.(bool)
}

func GetIntFromContext(ctx context.Context, key define.ContextKey) int {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	return value.(int)
}

func GetInt64FromContext(ctx context.Context, key define.ContextKey) int64 {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	return value.(int64)
}

func GetFloat64FromContext(ctx context.Context, key define.ContextKey) float64 {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	return value.(float64)
}

func GetInterfaceFromContext(ctx context.Context, key define.ContextKey) interface{} {
	value := ctx.Value(key)
	if value == nil {
		return nil
	}
	return value
}

func SetIntoContext(ctx context.Context, values map[define.ContextKey]interface{}) context.Context {
	for k, v := range values {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}
