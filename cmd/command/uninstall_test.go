package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Uninstall", func() {
	It("", func() {
		checkCommandFlag(uninstallCmd, "name", "n", cmdr.CfgKeyCommandUninstallName, "", true)
		checkCommandFlag(uninstallCmd, "version", "v", cmdr.CfgKeyCommandUninstallVersion, "", true)
	})
})
