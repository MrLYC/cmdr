package main

import (
	"os"
	"path/filepath"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/utils"
)

func init() {
	homeDir, err := os.UserHomeDir()
	utils.CheckError(err)

	cfg := config.Global

	cfg.SetDefault(config.CfgKeyCmdrRoot, filepath.Join(homeDir, ".cmdr"))
	cfg.SetDefault(config.CfgKeyBinDir, "bin")
	cfg.SetDefault(config.CfgKeyShimsDir, "shims")
	cfg.SetDefault(config.CfgKeyDatabase, "cmdr.db")
	cfg.SetDefault(config.CfgKeyLogLevel, "info")
	cfg.SetDefault(config.CfgKeyLogOutput, "stderr")
}
