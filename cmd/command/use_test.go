package command

import (
	"github.com/mrlyc/cmdr/config"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Use", func() {
	It("", func() {
		checkCommandFlag(useCmd, "name", "n", config.CfgKeyCommandUseName, "", true)
		checkCommandFlag(useCmd, "version", "v", config.CfgKeyCommandUseVersion, "", true)
	})
})
