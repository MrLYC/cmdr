package runner

import (
	"os"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandListName     = "command.list.name"
	CfgKeyCommandListVersion  = "command.list.version"
	CfgKeyCommandListLocation = "command.list.location"
	CfgKeyCommandListActivate = "command.list.activate"
)

func NewListRunner(cfg define.Configuration) *Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewFullCommandsQuerier(
			cfg.GetString(CfgKeyCommandInstallName),
			cfg.GetString(CfgKeyCommandInstallVersion),
			cfg.GetString(CfgKeyCommandInstallLocation),
			cfg.GetBool(CfgKeyCommandListActivate),
		),
		operator.NewCommandSorter(),
		operator.NewCommandPrinter(os.Stdout),
	)
}
