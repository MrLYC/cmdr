package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewInstallRunner(cfg define.Configuration) Runner {
	binDir := cfg.GetString(config.CfgKeyBinDir)
	shimsDir := cfg.GetString(config.CfgKeyShimsDir)

	runner := New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			shimsDir,
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			true,
		),
		operator.NewDownloader(),
		operator.NewBinariesInstaller(shimsDir),
	)

	if cfg.GetBool(config.CfgKeyCommandInstallActivate) {
		runner.Add(
			operator.NewCommandDeactivator(),
			operator.NewBinariesActivator(binDir, shimsDir),
			operator.NewCommandActivator(),
		)
	}

	return runner
}
