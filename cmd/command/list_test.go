package command

import (
	"github.com/mrlyc/cmdr/config"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("List", func() {
	It("", func() {
		checkCommandFlag(listCmd, "name", "n", config.CfgKeyCommandListName, "", false)
		checkCommandFlag(listCmd, "version", "v", config.CfgKeyCommandListVersion, "", false)
		checkCommandFlag(listCmd, "location", "l", config.CfgKeyCommandListLocation, "", false)
		checkCommandFlag(listCmd, "activate", "a", config.CfgKeyCommandListActivate, "false", false)
	})
})
