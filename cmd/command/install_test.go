package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/config"
)

var _ = Describe("Install", func() {
	It("", func() {
		checkCommandFlag(installCmd, "name", "n", config.CfgKeyCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", config.CfgKeyCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", config.CfgKeyCommandInstallLocation, "", true)
	})
})
