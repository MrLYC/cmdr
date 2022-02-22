package core

import (
	"github.com/spf13/viper"
)

type Configuration = *viper.Viper

var globalConfiguration Configuration

func GetConfiguration() Configuration {
	return globalConfiguration
}

func SetConfiguration(cfg Configuration) {
	globalConfiguration = cfg
}

const (
	// cmdr
	CfgKeyCmdrRootDir      = "core.root_dir"
	CfgKeyCmdrBinDir       = "core.bin_dir"
	CfgKeyCmdrShimsDir     = "core.shims_dir"
	CfgKeyCmdrProfileDir   = "core.profile_dir"
	CfgKeyCmdrDatabasePath = "core.database_path"
	CfgKeyCmdrProfilePath  = "core.profile_path"
	CfgKeyCmdrShell        = "core.shell"
	CfgKeyCmdrConfigPath   = "core.config_path"

	// log
	CfgKeyLogLevel  = "log.level"
	CfgKeyLogOutput = "log.output"

	// command.define
	CfgKeyCommandDefineName     = "_command.define.name"
	CfgKeyCommandDefineVersion  = "_command.define.version"
	CfgKeyCommandDefineLocation = "_command.define.location"

	// command.install
	CfgKeyCommandInstallName     = "_command.install.name"
	CfgKeyCommandInstallVersion  = "_command.install.version"
	CfgKeyCommandInstallLocation = "_command.install.location"
	CfgKeyCommandInstallActivate = "_command.install.activate"

	// command.list
	CfgKeyCommandListName     = "_command.list.name"
	CfgKeyCommandListVersion  = "_command.list.version"
	CfgKeyCommandListLocation = "_command.list.location"
	CfgKeyCommandListActivate = "_command.list.activate"

	// command.uninstall
	CfgKeyCommandUninstallName    = "_command.uninstall.name"
	CfgKeyCommandUninstallVersion = "_command.uninstall.version"

	// command.unset
	CfgKeyCommandUnsetName = "_command.unset.name"

	// command.use
	CfgKeyCommandUseName    = "_command.use.name"
	CfgKeyCommandUseVersion = "_command.use.version"
)

func init() {
	SetConfiguration(viper.GetViper())
}
