package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// getCmd represents the config command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get configuration by key",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		key := cfg.GetString(core.CfgKeyXConfigGetKey)
		value := cfg.Get(key)

		fmt.Printf("key: %s, type: %T, value: %v\n", key, value, value)
	},
}

func init() {
	Cmd.AddCommand(getCmd)

	flags := getCmd.Flags()
	flags.StringP("key", "k", "", "configuration key")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXConfigGetKey, flags.Lookup("key")),
		getCmd.MarkFlagRequired("key"),
	)
}
