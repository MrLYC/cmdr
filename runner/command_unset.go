package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUnsetRunner(cfg define.Configuration, cmdr *core.Cmdr) define.Runner {
	return New(
		operator.NewNamedCommandsQuerier(cfg.GetString(config.CfgKeyCommandUnsetName)),
		operator.NewCommandsChecker(),
		operator.NewBinariesDeactivator(cmdr),
		operator.NewCommandDeactivator(cmdr),
	)
}
