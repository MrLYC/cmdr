package config

import (
	"os"
	"path/filepath"

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
		configFile := cfg.ConfigFileUsed()

		userCfg := core.NewConfiguration()
		userCfg.SetConfigFile(configFile)
		err := userCfg.ReadInConfig()
		if err != nil {
			logger.Warn("failed to read user configuration file", map[string]interface{}{
				"file": configFile,
			})

			configDir := filepath.Dir(configFile)
			utils.PanicOnError("create configuration dir", os.MkdirAll(configDir, 0644))
		}

		key := cfg.GetString(core.CfgKeyXConfigSetKey)
		value := cfg.GetString(core.CfgKeyXConfigSetValue)

		userCfg.Set(key, value)

		logger.Info("writing user configuration", map[string]interface{}{
			"file": configFile,
		})
		utils.PanicOnError("write user configuration", userCfg.WriteConfig())
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
