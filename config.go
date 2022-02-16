package main

import (
	"os"
	"path/filepath"

	"github.com/mrlyc/cmdr/core"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfg := core.GetConfiguration()

	cfg.SetDefault(core.CfgKeyCmdrBinDir, filepath.Join(homeDir, "bin"))
	cfg.SetDefault(core.CfgKeyCmdrShimsDir, filepath.Join(homeDir, "shims"))
	cfg.SetDefault(core.CfgKeyLogLevel, "info")
	cfg.SetDefault(core.CfgKeyLogOutput, "stderr")
}
