package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"logur.dev/logur"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

var (
	cfgFile  string
	exitCode int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmdr",
	Short: "CMDR is a version manager for simple commands",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteContext(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

func init() {
	cobra.OnInitialize(preInitConfig, initConfig, postInitConfig, initLogger)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	pFlags := rootCmd.PersistentFlags()

	pFlags.StringVar(&cfgFile, "config", "config.yaml", "config file")

	cfg := core.GetConfiguration()
	utils.PanicOnError("binding flags", cfg.BindPFlag(core.CfgKeyCmdrConfigPath, pFlags.Lookup("config")))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preInitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfg := core.GetConfiguration()

	cfg.SetDefault(core.CfgKeyCmdrRootDir, filepath.Join(homeDir, ".cmdr"))
	cfg.SetDefault(core.CfgKeyCmdrBinDir, "bin")
	cfg.SetDefault(core.CfgKeyCmdrShimsDir, "shims")
	cfg.SetDefault(core.CfgKeyCmdrProfileDir, "profile")
	cfg.SetDefault(core.CfgKeyCmdrDatabasePath, "cmdr.db")
	cfg.SetDefault(core.CfgKeyCmdrShell, os.Getenv("SHELL"))

	cfg.SetDefault(core.CfgKeyLogLevel, "info")
	cfg.SetDefault(core.CfgKeyLogOutput, "stderr")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg := core.GetConfiguration()

	_, err := os.Stat(cfgFile)

	if err == nil {
		// Use config file from the flag.
		cfg.SetConfigFile(cfgFile)
		utils.CheckError(cfg.ReadInConfig())
	}

	cfg.AutomaticEnv() // read in environment variables that match
}

func postInitConfig() {
	cfg := core.GetConfiguration()
	rootDir := cfg.GetString(core.CfgKeyCmdrRootDir)

	for _, key := range []string{
		core.CfgKeyCmdrBinDir,
		core.CfgKeyCmdrShimsDir,
		core.CfgKeyCmdrProfileDir,
		core.CfgKeyCmdrDatabasePath,
		core.CfgKeyCmdrConfigPath,
	} {
		path := cfg.GetString(key)
		if !filepath.IsAbs(path) {
			cfg.Set(key, filepath.Join(rootDir, path))
		}
	}
}

func initLogger() {
	cfg := core.GetConfiguration()
	level, ok := logur.ParseLevel(cfg.GetString(core.CfgKeyLogLevel))
	if !ok {
		level = logur.Info
	}

	switch cfg.GetString(core.CfgKeyLogOutput) {
	case "stdout":
		color.SetOutput(os.Stdout)
	default:
		color.SetOutput(os.Stderr)
	}

	core.InitTerminalLogger(level, level < logur.Info, "error")
}
