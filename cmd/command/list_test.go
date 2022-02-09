package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("List", func() {
	It("", func() {
		checkCommandFlag(listCmd, "name", "n", cmdr.CfgKeyCommandListName, "", false)
		checkCommandFlag(listCmd, "version", "v", cmdr.CfgKeyCommandListVersion, "", false)
		checkCommandFlag(listCmd, "location", "l", cmdr.CfgKeyCommandListLocation, "", false)
		checkCommandFlag(listCmd, "activate", "a", cmdr.CfgKeyCommandListActivate, "false", false)
	})
})
