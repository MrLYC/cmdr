package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Use", func() {
	It("", func() {
		checkCommandFlag(useCmd, "name", "n", core.CfgKeyCommandUseName, "", true)
		checkCommandFlag(useCmd, "version", "v", core.CfgKeyCommandUseVersion, "", true)
	})
})
