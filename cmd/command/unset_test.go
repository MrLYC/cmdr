package command

import (
	"github.com/mrlyc/cmdr/config"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Unset", func() {
	It("", func() {
		checkCommandFlag(unsetCmd, "name", "n", config.CfgKeyCommandUnsetName, "", true)
	})
})
