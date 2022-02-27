package cmd

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Init", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(initCmd, "upgrade", "u", core.CfgKeyXInitUpgrade, "false", false)
	})
})
