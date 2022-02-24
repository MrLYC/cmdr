package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Define", func() {
	It("", func() {
		checkCommandFlag(defineCmd, "name", "n", core.CfgKeyXCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", core.CfgKeyXCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", core.CfgKeyXCommandDefineLocation, "", true)
	})
})
