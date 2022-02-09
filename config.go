package main

import (
	"os"
	"path/filepath"

	"github.com/mrlyc/cmdr/cmdr"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfg := cmdr.GetConfiguration()

	cfg.SetDefault(cmdr.CfgKeyCmdrRoot, filepath.Join(homeDir, ".cmdr"))
	cfg.SetDefault(cmdr.CfgKeyCmdrBinDir, filepath.Join(homeDir, "bin"))
	cfg.SetDefault(cmdr.CfgKeyCmdrShimsDir, filepath.Join(homeDir, "shims"))
	cfg.SetDefault(cmdr.CfgKeyLogLevel, "info")
	cfg.SetDefault(cmdr.CfgKeyLogOutput, "stderr")
}
