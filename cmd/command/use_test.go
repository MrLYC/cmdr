package command

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Use", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(UseCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandUseName, IsRequired: true},
			testutils.CommandFlagSpec{Name: "version", Shorthand: "v", ConfigKey: core.CfgKeyXCommandUseVersion, IsRequired: true},
		)
	})

	Context("command", func() {
		var (
			harness *testutils.CommandHarness
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDefault)

			harness.Cfg.Set(core.CfgKeyXCommandUseName, "cmdr")
			harness.Cfg.Set(core.CfgKeyXCommandUseVersion, "1.0.0")
		})

		AfterEach(func() {
			harness.Finish()
		})

		It("should unset a command", func() {
			harness.Manager.EXPECT().Activate("cmdr", "1.0.0").Return(nil)
			harness.ExpectClose()

			UseCmd.Run(UnsetCmd, []string{})
		})

		It("should panic when activate fails", func() {
			harness.Manager.EXPECT().Activate("cmdr", "1.0.0").Return(fmt.Errorf("activate failed"))
			harness.ExpectClose()

			Expect(func() { UseCmd.Run(UseCmd, []string{}) }).To(Panic())
		})
	})
})
