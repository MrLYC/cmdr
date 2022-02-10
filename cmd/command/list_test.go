package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("List", func() {
	It("", func() {
		checkCommandFlag(listCmd, "name", "n", core.CfgKeyCommandListName, "", false)
		checkCommandFlag(listCmd, "version", "v", core.CfgKeyCommandListVersion, "", false)
		checkCommandFlag(listCmd, "location", "l", core.CfgKeyCommandListLocation, "", false)
		checkCommandFlag(listCmd, "activate", "a", core.CfgKeyCommandListActivate, "false", false)
	})
})
