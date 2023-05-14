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

func NewConfiguration() *viper.Viper {
	return viper.New()
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
	CfgKeyCmdrLinkMode     = "core.link_mode"

	// proxy
	CfgKeyProxyGo    = "proxy.go"
	CfgKeyProxyHTTP  = "proxy.http"
	CfgKeyProxyHTTPS = "proxy.https"

	// log
	CfgKeyLogLevel  = "log.level"
	CfgKeyLogOutput = "log.output"

	// cmd.command.define
	CfgKeyXCommandDefineName     = "_.command.define.name"
	CfgKeyXCommandDefineVersion  = "_.command.define.version"
	CfgKeyXCommandDefineLocation = "_.command.define.location"
	CfgKeyXCommandDefineActivate = "_.command.define.activate"
	// cmd.command.install
	CfgKeyXCommandInstallName     = "_.command.install.name"
	CfgKeyXCommandInstallVersion  = "_.command.install.version"
	CfgKeyXCommandInstallLocation = "_.command.install.location"
	CfgKeyXCommandInstallActivate = "_.command.install.activate"
	// cmd.command.list
	CfgKeyXCommandListName     = "_.command.list.name"
	CfgKeyXCommandListVersion  = "_.command.list.version"
	CfgKeyXCommandListLocation = "_.command.list.location"
	CfgKeyXCommandListActivate = "_.command.list.activate"
	CfgKeyXCommandListFields   = "_.command.list.fields"
	// cmd.command.remove
	CfgKeyXCommandRemoveName    = "_.command.remove.name"
	CfgKeyXCommandRemoveVersion = "_.command.remove.version"
	// cmd.command.unset
	CfgKeyXCommandUnsetName = "_.command.unset.name"
	// cmd.command.use
	CfgKeyXCommandUseName    = "_.command.use.name"
	CfgKeyXCommandUseVersion = "_.command.use.version"

	// cmd.config.get
	CfgKeyXConfigGetKey = "_.config.get.key"
	// cmd.config.set
	CfgKeyXConfigSetKey   = "_.config.set.key"
	CfgKeyXConfigSetValue = "_.config.set.value"

	// cmd.init
	CfgKeyXInitUpgrade = "_.init.upgrade"

	// cmd.upgrade
	CfgKeyXUpgradeRelease = "_.upgrade.release"
	CfgKeyXUpgradeAsset   = "_.upgrade.asset"
	CfgKeyXUpgradeArgs    = "_.upgrade.args"
)

func init() {
	SetConfiguration(NewConfiguration())
}
