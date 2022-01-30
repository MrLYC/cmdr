package runner

import (
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandUseName    = "command.use.name"
	CfgKeyCommandUseVersion = "command.use.version"
)

func NewUseRunner(cfg define.Configuration) *Runner {
	binDir := operator.GetBinDir()
	shimsDir := operator.GetShimsDir()

	return New(
		operator.NewDBClientMaker(),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(CfgKeyCommandUseName),
			cfg.GetString(CfgKeyCommandUseVersion),
		).StrictMode(),
		operator.NewCommandDeactivator(),
		operator.NewBinariesActivator(binDir, shimsDir),
		operator.NewCommandActivator(),
	)
}
