package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("Define", func() {
	It("", func() {
		checkCommandFlag(defineCmd, "name", "n", runner.CfgKeyCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", runner.CfgKeyCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", runner.CfgKeyCommandDefineLocation, "", true)
	})
})
