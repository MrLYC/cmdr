package command

import (
	"github.com/mrlyc/cmdr/config"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Uninstall", func() {
	It("", func() {
		checkCommandFlag(uninstallCmd, "name", "n", config.CfgKeyCommandUninstallName, "", true)
		checkCommandFlag(uninstallCmd, "version", "v", config.CfgKeyCommandUninstallVersion, "", true)
	})
})
