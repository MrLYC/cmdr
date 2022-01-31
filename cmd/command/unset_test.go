package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("Unset", func() {
	It("", func() {
		checkCommandFlag(unsetCmd, "name", "n", runner.CfgKeyCommandUnsetName, "", true)
	})
})
