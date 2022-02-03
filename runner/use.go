package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewUseRunner(cfg define.Configuration, helper *utils.CmdrHelper) define.Runner {
	return New(
		operator.NewDBClientMaker(helper),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUseName),
			cfg.GetString(config.CfgKeyCommandUseVersion),
		).StrictMode(),
		operator.NewCommandDeactivator(),
		operator.NewBinariesActivator(helper),
		operator.NewCommandActivator(),
	)
}
