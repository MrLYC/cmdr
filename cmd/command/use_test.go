package command

import (
	"github.com/mrlyc/cmdr/cmdr"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Use", func() {
	It("", func() {
		checkCommandFlag(useCmd, "name", "n", cmdr.CfgKeyCommandUseName, "", true)
		checkCommandFlag(useCmd, "version", "v", cmdr.CfgKeyCommandUseVersion, "", true)
	})
})
