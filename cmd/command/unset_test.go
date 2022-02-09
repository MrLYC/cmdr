package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Unset", func() {
	It("", func() {
		checkCommandFlag(unsetCmd, "name", "n", cmdr.CfgKeyCommandUnsetName, "", true)
	})
})
