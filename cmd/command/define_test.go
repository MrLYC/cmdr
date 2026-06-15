package command

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Define", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(DefineCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandDefineName, IsRequired: true},
			testutils.CommandFlagSpec{Name: "version", Shorthand: "v", ConfigKey: core.CfgKeyXCommandDefineVersion, IsRequired: true},
			testutils.CommandFlagSpec{Name: "location", Shorthand: "l", ConfigKey: core.CfgKeyXCommandDefineLocation, IsRequired: true},
			testutils.CommandFlagSpec{Name: "activate", Shorthand: "a", ConfigKey: core.CfgKeyXCommandDefineActivate, Default: "false"},
		)
	})

	Context("command", func() {
		var (
			harness *testutils.CommandHarness
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDefault)

			harness.Cfg.Set(core.CfgKeyXCommandDefineName, "test")
			harness.Cfg.Set(core.CfgKeyXCommandDefineVersion, "1.0.0")
			harness.Cfg.Set(core.CfgKeyXCommandDefineLocation, "")
		})

		AfterEach(func() {
			harness.Finish()
		})

		DescribeTable("should define commands",
			func(activate bool) {
				harness.Cfg.Set(core.CfgKeyXCommandDefineActivate, activate)

				harness.Manager.EXPECT().Define("test", "1.0.0", "")
				if activate {
					harness.Manager.EXPECT().Activate("test", "1.0.0").Return(nil)
				}
				harness.ExpectClose()

				DefineCmd.Run(DefineCmd, []string{})
			},
			Entry("and activates it", true),
			Entry("without activation", false),
		)

		It("should panic when define fails", func() {
			harness.Cfg.Set(core.CfgKeyXCommandDefineActivate, false)
			harness.Manager.EXPECT().Define("test", "1.0.0", "").Return(nil, fmt.Errorf("define failed"))
			harness.ExpectClose()

			Expect(func() { DefineCmd.Run(DefineCmd, []string{}) }).To(Panic())
		})

		It("should change link mode", func() {
			DefineCmd.PreRun(DefineCmd, []string{})

			Expect(harness.Cfg.GetString(core.CfgKeyCmdrLinkMode)).To(Equal("link"))
		})
	})
})
