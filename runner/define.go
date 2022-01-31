package runner

import (
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandDefineName     = "command.define.name"
	CfgKeyCommandDefineVersion  = "command.define.version"
	CfgKeyCommandDefineLocation = "command.define.location"
)

func NewDefineRunner(cfg define.Configuration) Runner {
	shimsDir := operator.GetShimsDir()

	return New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			shimsDir,
			cfg.GetString(CfgKeyCommandInstallName),
			cfg.GetString(CfgKeyCommandInstallVersion),
			cfg.GetString(CfgKeyCommandInstallLocation),
			false,
		),
	)
}
