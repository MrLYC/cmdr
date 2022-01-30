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

	cfg := define.Config

	cfg.SetDefault(define.CfgKeyCmdrRoot, filepath.Join(homeDir, ".cmdr"))
	cfg.SetDefault(define.CfgKeyBinDir, "bin")
	cfg.SetDefault(define.CfgKeyShimsDir, "shims")
	cfg.SetDefault(define.CfgKeyDatabase, "cmdr.db")
	cfg.SetDefault(define.CfgKeyLogLevel, "info")
	cfg.SetDefault(define.CfgKeyLogOutput, "stderr")
}
