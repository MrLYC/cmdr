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
		out, err := renderConfigValue(cfg, key)
		utils.PanicOnError("render config failed", err)
		fmt.Printf("%s", out)
	},
}

func renderConfigValue(cfg core.Configuration, key string) ([]byte, error) {
	var dumps interface{}
	if err := cfg.UnmarshalKey(key, &dumps); err != nil {
		return nil, err
	}

	return yaml.Marshal(dumps)
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
