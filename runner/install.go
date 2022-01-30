package runner

import (
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandInstallName     = "command.install.name"
	CfgKeyCommandInstallVersion  = "command.install.version"
	CfgKeyCommandInstallLocation = "command.install.location"
	CfgKeyCommandInstallActivate = "command.install.activate"
)

func NewInstallRunner(cfg define.Configuration) *Runner {
	binDir := cfg.GetString(define.CfgKeyBinDir)
	shimsDir := cfg.GetString(define.CfgKeyShimsDir)

	runner := New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			shimsDir,
			cfg.GetString(CfgKeyCommandInstallName),
			cfg.GetString(CfgKeyCommandInstallVersion),
			cfg.GetString(CfgKeyCommandInstallLocation),
			true,
		),
		operator.NewDownloader(),
		operator.NewBinariesInstaller(shimsDir),
	)

	if cfg.GetBool(CfgKeyCommandInstallActivate) {
		runner.Add(
			operator.NewCommandDeactivator(),
			operator.NewBinariesActivator(binDir, shimsDir),
			operator.NewCommandActivator(),
		)
	}

	return runner
}
