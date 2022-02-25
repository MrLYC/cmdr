package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Remove", func() {
	It("should check flags", func() {
		checkCommandFlag(removeCmd, "name", "n", core.CfgKeyXCommandRemoveName, "", true)
		checkCommandFlag(removeCmd, "version", "v", core.CfgKeyXCommandRemoveVersion, "", true)
	})
})
