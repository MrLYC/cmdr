package config

import (
	"context"

	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/define"
)

var Global define.Configuration

func GetGlobalConfiguration() define.Configuration {
	return Global
}

func GetConfigurationFromContext(ctx context.Context) define.Configuration {
	cfg := ctx.Value(define.ContextKeyConfiguration)
	if cfg == nil {
		return Global
	}
	return cfg.(define.Configuration)
}

func init() {
	Global = viper.GetViper()
}
