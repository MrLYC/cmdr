package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Install", func() {
	It("", func() {
		checkCommandFlag(installCmd, "name", "n", cmdr.CfgKeyCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", cmdr.CfgKeyCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", cmdr.CfgKeyCommandInstallLocation, "", true)
	})
})
