package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUseRunner(cfg define.Configuration, cmdr *core.Cmdr) define.Runner {
	return New(
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUseName),
			cfg.GetString(config.CfgKeyCommandUseVersion),
		),
		operator.NewCommandsChecker(),
		operator.NewCommandDeactivator(cmdr),
		operator.NewBinariesActivator(cmdr),
		operator.NewCommandActivator(cmdr),
	)
}
