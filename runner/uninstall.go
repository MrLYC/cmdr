package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUninstallRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUninstallName),
			cfg.GetString(config.CfgKeyCommandUninstallVersion),
		).StrictMode(),
		operator.NewCommandUndefiner(),
		operator.NewBinariesUninstaller(),
	)
}
