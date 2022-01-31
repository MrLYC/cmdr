package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("List", func() {
	It("", func() {
		checkCommandFlag(listCmd, "name", "n", runner.CfgKeyCommandListName, "", false)
		checkCommandFlag(listCmd, "version", "v", runner.CfgKeyCommandListVersion, "", false)
		checkCommandFlag(listCmd, "location", "l", runner.CfgKeyCommandListLocation, "", false)
		checkCommandFlag(listCmd, "activate", "a", runner.CfgKeyCommandListActivate, "false", false)
	})
})
