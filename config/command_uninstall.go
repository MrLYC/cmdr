package config

const (
	// command.uninstall
	CfgKeyCommandUninstallName    = "command.uninstall.name"
	CfgKeyCommandUninstallVersion = "command.uninstall.version"
)

func GetCommandUninstallName() string {
	return CfgKeyCommandUninstallName
}

func GetCommandUninstallVersion() string {
	return CfgKeyCommandUninstallVersion
}
