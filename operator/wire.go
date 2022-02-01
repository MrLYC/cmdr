// +build wireinject

package operator

import (
	"context"

	"github.com/google/wire"
	"github.com/mrlyc/cmdr/config"
)

var (
	binDirSet   = wire.NewSet(config.GetConfigurationFromContext, config.GetBinDir)
	shimsDirSet = wire.NewSet(config.GetConfigurationFromContext, config.GetShimsDir)
	dbSet       = wire.NewSet(config.GetConfigurationFromContext, config.GetDatabasePath)
)

func GetBinDir(ctx context.Context) string {
	wire.Build(binDirSet)
	return ""
}

func GetShimsDir(ctx context.Context) string {
	wire.Build(binDirSet)
	return ""
}

func GetDatabasePath(ctx context.Context) string {
	wire.Build(dbSet)
	return ""
}
