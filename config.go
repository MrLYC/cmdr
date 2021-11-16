package main

import (
	"os"
	"path/filepath"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

func init() {
	homeDir, err := os.UserHomeDir()
	utils.CheckError(err)

	cfg := define.Configuration

	cfg.SetDefault("cmdr.root", filepath.Join(homeDir, ".cmdr"))

	cfg.SetDefault("database.name", "cmdr.db")
	cfg.SetDefault("log.level", "info")
	cfg.SetDefault("log.output", "stderr")
}
