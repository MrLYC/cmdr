package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Unset", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(UnsetCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandUnsetName, IsRequired: true},
		)
	})

	Context("command", func() {
		var (
			harness *testutils.CommandHarness
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDefault)

			harness.Cfg.Set(core.CfgKeyXCommandUnsetName, "testing")
		})

		AfterEach(func() {
			harness.Finish()
		})

		It("should unset a command", func() {
			harness.Manager.EXPECT().Deactivate("testing").Return(nil)
			harness.ExpectClose()

			UnsetCmd.Run(UnsetCmd, []string{})
		})

		It("should not unset a cmdr", func() {
			harness.Cfg.Set(core.CfgKeyXCommandUnsetName, "cmdr")
			harness.ExpectClose()

			UnsetCmd.Run(UnsetCmd, []string{})
		})
	})
})
