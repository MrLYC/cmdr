package config

const (
	// command.use
	CfgKeyCommandUseName    = "command.use.name"
	CfgKeyCommandUseVersion = "command.use.version"
)

func GetCommandUseName() string {
	return CfgKeyCommandUseName
}

func GetCommandUseVersion() string {
	return CfgKeyCommandUseVersion
}
