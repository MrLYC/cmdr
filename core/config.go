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
	CfgKeyCmdrBinDir       = "core.bin_dir"
	CfgKeyCmdrShimsDir     = "core.shims_dir"
	CfgKeyCmdrProfileDir   = "core.profile_dir"
	CfgKeyCmdrDatabasePath = "core.database_path"
	CfgKeyCmdrProfilePath  = "core.profile_path"
	CfgKeyCmdrProfileName  = "core.profile_name"

	// log
	CfgKeyLogLevel  = "log.level"
	CfgKeyLogOutput = "log.output"

	// command.define
	CfgKeyCommandDefineName     = "command.define.name"
	CfgKeyCommandDefineVersion  = "command.define.version"
	CfgKeyCommandDefineLocation = "command.define.location"

	// command.install
	CfgKeyCommandInstallName     = "command.install.name"
	CfgKeyCommandInstallVersion  = "command.install.version"
	CfgKeyCommandInstallLocation = "command.install.location"
	CfgKeyCommandInstallActivate = "command.install.activate"

	// command.list
	CfgKeyCommandListName     = "command.list.name"
	CfgKeyCommandListVersion  = "command.list.version"
	CfgKeyCommandListLocation = "command.list.location"
	CfgKeyCommandListActivate = "command.list.activate"

	// command.uninstall
	CfgKeyCommandUninstallName    = "command.uninstall.name"
	CfgKeyCommandUninstallVersion = "command.uninstall.version"

	// command.unset
	CfgKeyCommandUnsetName = "command.unset.name"

	// command.use
	CfgKeyCommandUseName    = "command.use.name"
	CfgKeyCommandUseVersion = "command.use.version"
)

func init() {
	SetConfiguration(viper.GetViper())
}
