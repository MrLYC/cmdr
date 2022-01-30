package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

var installCmdFlag struct {
	activate bool
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewInstallRunner(define.Config)

		utils.ExitWithError(runner.Run(cmd.Context()), "install failed")

		define.Logger.Info("installed command", map[string]interface{}{
			"name":     simpleCmdFlag.name,
			"version":  simpleCmdFlag.version,
			"location": simpleCmdFlag.location,
		})
	},
}

func init() {
	Cmd.AddCommand(installCmd)
	cmdFlagsHelper.declareFlagName(installCmd)
	cmdFlagsHelper.declareFlagVersion(installCmd)

	flags := installCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.location, "location", "l", "", "command location")
	flags.BoolVarP(&installCmdFlag.activate, "activate", "a", false, "activate command")

	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandInstallName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandInstallVersion, flags.Lookup("version"))
	cfg.BindPFlag(runner.CfgKeyCommandInstallLocation, flags.Lookup("location"))
	cfg.BindPFlag(runner.CfgKeyCommandInstallActivate, flags.Lookup("activate"))

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
