package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Install", func() {
	It("", func() {
		checkCommandFlag(installCmd, "name", "n", core.CfgKeyCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", core.CfgKeyCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", core.CfgKeyCommandInstallLocation, "", true)
	})
})
