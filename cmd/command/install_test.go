package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("Install", func() {
	It("", func() {
		checkCommandFlag(installCmd, "name", "n", runner.CfgKeyCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", runner.CfgKeyCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", runner.CfgKeyCommandInstallLocation, "", true)
	})
})
