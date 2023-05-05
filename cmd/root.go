package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"logur.dev/logur"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
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
		utils.ExitOnError("execute failed", err)
	}
}

func init() {
	cobra.OnInitialize(preInitConfig, initConfig, postInitConfig, initLogger, initDatabase)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global forgo get -u github.com/ory/dockertest/v3 your application.

	cfg := core.GetConfiguration()
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringP(
		"config", "c",
		filepath.Join(cfg.GetString(core.CfgKeyCmdrRootDir), "config.yaml"),
		"config file",
	)

	utils.PanicOnError("binding flags", cfg.BindPFlag(core.CfgKeyCmdrConfigPath, pFlags.Lookup("config")))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preInitConfig() {
	cfg := core.GetConfiguration()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.SetEnvPrefix("cmdr")
	cfg.AutomaticEnv() // read in environment variables that match

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfg.SetDefault(core.CfgKeyCmdrRootDir, filepath.Join(homeDir, ".cmdr"))
	cfg.SetDefault(core.CfgKeyCmdrBinDir, "bin")
	cfg.SetDefault(core.CfgKeyCmdrShimsDir, "shims")
	cfg.SetDefault(core.CfgKeyCmdrProfileDir, "profile")
	cfg.SetDefault(core.CfgKeyCmdrDatabasePath, "cmdr.db")

	cfg.SetDefault(core.CfgKeyLogLevel, "info")
	cfg.SetDefault(core.CfgKeyLogOutput, "stderr")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg := core.GetConfiguration()

	cfgFile := cfg.GetString(core.CfgKeyCmdrConfigPath)

	cfg.SetConfigFile(cfgFile)
	_, err := os.Stat(cfgFile)
	if err == nil {
		// Use config file from the flag.
		utils.CheckError(cfg.ReadInConfig())
	}
}

func postInitConfig() {
	cfg := core.GetConfiguration()
	rootDir := cfg.GetString(core.CfgKeyCmdrRootDir)

	for _, key := range []string{
		core.CfgKeyCmdrBinDir,
		core.CfgKeyCmdrShimsDir,
		core.CfgKeyCmdrProfileDir,
		core.CfgKeyCmdrDatabasePath,
	} {
		path := cfg.GetString(key)
		if filepath.IsAbs(path) {
			continue
		}

		value := filepath.Join(rootDir, path)

		if cfg.IsSet(key) {
			cfg.SetDefault(key, value)
		} else {
			cfg.SetDefault(key, value)
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

func initDatabase() {
	core.SetDatabaseFactory(func() (core.Database, error) {
		logger := core.GetLogger()

		core.SubscribeEventOnce(core.EventExit, func() {
			utils.PanicOnError("closing database", core.CloseDatabase())
		})

		cfg := core.GetConfiguration()
		dbPath := cfg.GetString(core.CfgKeyCmdrDatabasePath)
		logger.Debug("opening database", map[string]interface{}{
			"path": dbPath,
		})

		db, err := storm.Open(dbPath)
		if err != nil {
			return nil, errors.Wrapf(err, "open database failed")
		}

		return db, nil
	})
}
