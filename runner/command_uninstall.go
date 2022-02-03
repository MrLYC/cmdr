package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewUninstallRunner(cfg define.Configuration, helper *utils.CmdrHelper) define.Runner {
	return New(
		operator.NewDBClientMaker(helper),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUninstallName),
			cfg.GetString(config.CfgKeyCommandUninstallVersion),
		).StrictMode(),
		operator.NewCommandUndefiner(),
		operator.NewBinariesUninstaller(),
	)
}
