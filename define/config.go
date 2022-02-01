package define

import "github.com/spf13/viper"

type Configuration = *viper.Viper

const (
	CfgKeyCmdrRoot = "cmdr.root_dir"
	CfgKeyBinDir   = "cmdr.binary_dir"
	CfgKeyCmdDir   = "cmdr.command_dir"
	CfgKeyShimsDir = "cmdr.shims_dir"
	CfgKeyDatabase = "cmdr.database"
)
