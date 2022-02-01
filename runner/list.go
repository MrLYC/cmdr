package runner

import (
	"os"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewListRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewFullCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			cfg.GetBool(config.CfgKeyCommandListActivate),
		),
		operator.NewCommandSorter(),
		operator.NewCommandPrinter(os.Stdout),
	)
}
