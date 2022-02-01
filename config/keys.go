package config

const (
	// cmdr
	CfgKeyCmdrRoot = "cmdr.root_dir"
	CfgKeyBinDir   = "cmdr.binary_dir"
	CfgKeyCmdDir   = "cmdr.command_dir"
	CfgKeyShimsDir = "cmdr.shims_dir"
	CfgKeyDatabase = "cmdr.database"

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

	// command.update
	CfgKeyCommandUnsetName = "command.unset.name"

	// command.use
	CfgKeyCommandUseName    = "command.use.name"
	CfgKeyCommandUseVersion = "command.use.version"
)
