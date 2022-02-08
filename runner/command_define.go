package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewDefineRunner(cfg define.Configuration, cmdr *core.Cmdr) define.Runner {
	return New(
		operator.NewCommandDefiner(
			cmdr,
			cfg.GetString(config.CfgKeyCommandDefineName),
			cfg.GetString(config.CfgKeyCommandDefineVersion),
			cfg.GetString(config.CfgKeyCommandDefineLocation),
		),
		operator.NewBinariesChecker(),
	)
}
