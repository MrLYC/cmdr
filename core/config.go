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
	CfgKeyXCommandDefineName     = "_.command.define.name"
	CfgKeyXCommandDefineVersion  = "_.command.define.version"
	CfgKeyXCommandDefineLocation = "_.command.define.location"
	// command.install
	CfgKeyXCommandInstallName     = "_.command.install.name"
	CfgKeyXCommandInstallVersion  = "_.command.install.version"
	CfgKeyXCommandInstallLocation = "_.command.install.location"
	CfgKeyXCommandInstallActivate = "_.command.install.activate"
	// command.list
	CfgKeyXCommandListName     = "_.command.list.name"
	CfgKeyXCommandListVersion  = "_.command.list.version"
	CfgKeyXCommandListLocation = "_.command.list.location"
	CfgKeyXCommandListActivate = "_.command.list.activate"
	// command.uninstall
	CfgKeyXCommandUninstallName    = "_.command.uninstall.name"
	CfgKeyXCommandUninstallVersion = "_.command.uninstall.version"
	// command.unset
	CfgKeyXCommandUnsetName = "_.command.unset.name"
	// command.use
	CfgKeyXCommandUseName    = "_.command.use.name"
	CfgKeyXCommandUseVersion = "_.command.use.version"

	// config.get
	CfgKeyXConfigGetKey = "_.config.get.key"
)

func init() {
	SetConfiguration(viper.GetViper())
}
