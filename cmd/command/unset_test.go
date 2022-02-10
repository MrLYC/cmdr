package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Unset", func() {
	It("", func() {
		checkCommandFlag(unsetCmd, "name", "n", core.CfgKeyCommandUnsetName, "", true)
	})
})
