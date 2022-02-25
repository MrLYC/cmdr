package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("List", func() {
	It("should check flags", func() {
		checkCommandFlag(listCmd, "name", "n", core.CfgKeyXCommandListName, "", false)
		checkCommandFlag(listCmd, "version", "v", core.CfgKeyXCommandListVersion, "", false)
		checkCommandFlag(listCmd, "location", "l", core.CfgKeyXCommandListLocation, "", false)
		checkCommandFlag(listCmd, "activate", "a", core.CfgKeyXCommandListActivate, "false", false)
	})
})
