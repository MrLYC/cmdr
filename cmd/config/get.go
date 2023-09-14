package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

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
		var dumps interface{}
		utils.PanicOnError("unmarshal config failed", cfg.UnmarshalKey(key, &dumps))

		out, err := yaml.Marshal(dumps)
		utils.PanicOnError("marshal config failed", err)
		fmt.Printf("%s", out)
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
