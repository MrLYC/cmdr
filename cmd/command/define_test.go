package command

import (
	"github.com/mrlyc/cmdr/core"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Define", func() {
	It("", func() {
		checkCommandFlag(defineCmd, "name", "n", core.CfgKeyCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", core.CfgKeyCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", core.CfgKeyCommandDefineLocation, "", true)
	})
})
