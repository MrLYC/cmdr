package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
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
	cobra.OnInitialize(initConfig, define.InitLogger)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	pFlags := rootCmd.PersistentFlags()

	homeDir, err := os.UserHomeDir()
	utils.CheckError(err)
	filepath.Join()
	pFlags.StringVar(&cfgFile, "config", filepath.Join(homeDir, ".cmdr.yaml"), "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg := define.Configuration

	_, err := define.FS.Stat(cfgFile)

	if err == nil {
		// Use config file from the flag.
		cfg.SetConfigFile(cfgFile)
		utils.CheckError(cfg.ReadInConfig())
	}

	cfg.AutomaticEnv() // read in environment variables that match
}
