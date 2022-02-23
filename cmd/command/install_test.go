package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Install", func() {
	It("", func() {
		checkCommandFlag(installCmd, "name", "n", core.CfgKeyXCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", core.CfgKeyXCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", core.CfgKeyXCommandInstallLocation, "", true)
	})
})
