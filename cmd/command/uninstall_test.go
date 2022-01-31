package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("Uninstall", func() {
	It("", func() {
		checkCommandFlag(uninstallCmd, "name", "n", runner.CfgKeyCommandUninstallName, "", true)
		checkCommandFlag(uninstallCmd, "version", "v", runner.CfgKeyCommandUninstallVersion, "", true)
	})
})
