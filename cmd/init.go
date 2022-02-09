//+build !windows

package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
	"github.com/mrlyc/cmdr/cmdr/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial cmdr environment",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmdr.GetConfiguration()

		tmpl, err := template.New("init.sh").ParseFS(cmdr.EmbedFS, "scripts/init.sh")
		utils.CheckError(err)

		var buffer bytes.Buffer
		utils.CheckError(tmpl.Execute(&buffer, map[string]interface{}{
			"BinDir": cfg.GetString(cmdr.CfgKeyCmdrDatabasePath),
		}))

		fmt.Println(buffer.String())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
