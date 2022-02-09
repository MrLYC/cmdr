package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Define", func() {
	It("", func() {
		checkCommandFlag(defineCmd, "name", "n", cmdr.CfgKeyCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", cmdr.CfgKeyCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", cmdr.CfgKeyCommandDefineLocation, "", true)
	})
})
