//+build !windows

package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial cmdr environment",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetGlobalConfiguration()
		tmpl, err := template.New("init.sh").ParseFS(define.EmbedFS, "scripts/init.sh")
		utils.CheckError(err)

		var buffer bytes.Buffer
		utils.CheckError(tmpl.Execute(&buffer, map[string]interface{}{
			"BinDir": config.GetBinDir(cfg),
		}))

		fmt.Println(buffer.String())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
