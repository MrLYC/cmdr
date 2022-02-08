package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUninstallRunner(cfg define.Configuration, cmdr *core.Cmdr) define.Runner {
	return New(
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUninstallName),
			cfg.GetString(config.CfgKeyCommandUninstallVersion),
		),
		operator.NewCommandsChecker(),
		operator.NewCommandUndefiner(cmdr),
		operator.NewBinariesUninstaller(cmdr),
	)
}
