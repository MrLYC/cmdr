package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewDefineRunner(cfg define.Configuration, helper *utils.CmdrHelper) define.Runner {
	return New(
		operator.NewDBClientMaker(helper),
		operator.NewCommandDefiner(
			cfg.GetString(config.CfgKeyCommandDefineName),
			cfg.GetString(config.CfgKeyCommandDefineVersion),
			cfg.GetString(config.CfgKeyCommandDefineLocation),
			false,
			helper,
		),
		operator.NewBinariesChecker(),
	)
}
