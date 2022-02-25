package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Unset", func() {
	It("should check flags", func() {
		checkCommandFlag(unsetCmd, "name", "n", core.CfgKeyXCommandUnsetName, "", true)
	})
})
