package command

import (
	"github.com/mrlyc/cmdr/config"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Define", func() {
	It("", func() {
		checkCommandFlag(defineCmd, "name", "n", config.CfgKeyCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", config.CfgKeyCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", config.CfgKeyCommandDefineLocation, "", true)
	})
})
