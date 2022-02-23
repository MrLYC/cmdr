package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Uninstall", func() {
	It("", func() {
		checkCommandFlag(uninstallCmd, "name", "n", core.CfgKeyXCommandUninstallName, "", true)
		checkCommandFlag(uninstallCmd, "version", "v", core.CfgKeyXCommandUninstallVersion, "", true)
	})
})
