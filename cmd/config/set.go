package config

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// setCmd represents the config command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.GetLogger()

		cfg := core.GetConfiguration()
		key := cfg.GetString(core.CfgKeyXConfigSetKey)
		value := cfg.GetString(core.CfgKeyXConfigSetValue)

		cfg.Set(key, value)
		cfg.Set("_", nil)

		logger.Info("writing configuration", map[string]interface{}{
			"file": cfg.ConfigFileUsed(),
		})
		utils.PanicOnError("write config", cfg.WriteConfig())
	},
}

func init() {
	Cmd.AddCommand(setCmd)

	flags := setCmd.Flags()
	flags.StringP("key", "k", "", "configuration key")
	flags.StringP("value", "v", "", "configuration value")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXConfigSetKey, flags.Lookup("key")),
		cfg.BindPFlag(core.CfgKeyXConfigSetValue, flags.Lookup("value")),
		setCmd.MarkFlagRequired("key"),
		setCmd.MarkFlagRequired("value"),
	)
}
