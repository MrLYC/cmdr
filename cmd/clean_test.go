package cmd

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Clean", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(cleanCmd, "age", "", core.CfgKeyXCleanAgeDays, "100", false)
		testutils.CheckCommandFlag(cleanCmd, "keep", "", core.CfgKeyXCleanKeep, "3", false)
		testutils.CheckCommandFlag(cleanCmd, "name", "n", core.CfgKeyXCleanName, "", false)
	})
})
