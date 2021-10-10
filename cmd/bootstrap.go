package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootrap cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, err := template.New("bootstrap.sh").ParseFS(define.EmbedFS, "scripts/bootstrap.sh")
		utils.CheckError(err)

		var buffer bytes.Buffer
		utils.CheckError(tmpl.Execute(&buffer, map[string]interface{}{
			"BinDir": core.GetBinDir(),
		}))

		fmt.Println(buffer.String())
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}
