package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("Use", func() {
	It("", func() {
		checkCommandFlag(useCmd, "name", "n", runner.CfgKeyCommandUseName, "", true)
		checkCommandFlag(useCmd, "version", "v", runner.CfgKeyCommandUseVersion, "", true)
	})
})
