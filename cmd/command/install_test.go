package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Install", func() {
	It("should check flags", func() {
		checkCommandFlag(installCmd, "name", "n", core.CfgKeyXCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", core.CfgKeyXCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", core.CfgKeyXCommandInstallLocation, "", true)
		checkCommandFlag(installCmd, "activate", "a", core.CfgKeyXCommandInstallActivate, "false", false)
	})
})
