package runner

import (
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandUninstallName    = "command.uninstall.name"
	CfgKeyCommandUninstallVersion = "command.uninstall.version"
)

func NewUninstallRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(CfgKeyCommandUninstallName),
			cfg.GetString(CfgKeyCommandUninstallVersion),
		).StrictMode(),
		operator.NewCommandUndefiner(),
		operator.NewBinariesUninstaller(),
	)
}
