package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

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
		yamlValue := cfg.GetString(core.CfgKeyXConfigSetValue)

		utils.PanicOnError("write user configuration", setConfigValue(cfg.ConfigFileUsed(), key, yamlValue, logger))
	},
}

func setConfigValue(configFile, key, yamlValue string, logger interface {
	Warn(string, ...map[string]interface{})
	Info(string, ...map[string]interface{})
}) error {
	userCfg := core.NewConfiguration()
	userCfg.SetConfigFile(configFile)
	err := userCfg.ReadInConfig()
	if err != nil {
		logger.Warn("failed to read user configuration file", map[string]interface{}{
			"file": configFile,
		})

		configDir := filepath.Dir(configFile)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	var value interface{}
	if err := yaml.Unmarshal([]byte(yamlValue), &value); err != nil {
		return err
	}

	userCfg.Set(key, value)

	logger.Info("writing user configuration", map[string]interface{}{
		"file": configFile,
	})
	return userCfg.WriteConfig()
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
