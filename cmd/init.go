package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial cmdr environment",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger

		cfg := define.Configuration
		cmdrDir := cfg.GetString("cmdr_dir")

		logger.Debug("Creating cmdr dir", map[string]interface{}{
			"dir": cmdrDir,
		})
		utils.CheckError(os.MkdirAll(cmdrDir, 0755))

		client := core.GetClient()
		defer utils.CallClose(client)

		logger.Debug("Creating cmdr database")
		utils.CheckError(client.Schema.Create(cmd.Context()))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
