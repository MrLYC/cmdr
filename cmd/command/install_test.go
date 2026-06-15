package command

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Install", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(InstallCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandInstallName, IsRequired: true},
			testutils.CommandFlagSpec{Name: "version", Shorthand: "v", ConfigKey: core.CfgKeyXCommandInstallVersion, IsRequired: true},
			testutils.CommandFlagSpec{Name: "location", Shorthand: "l", ConfigKey: core.CfgKeyXCommandInstallLocation, IsRequired: true},
			testutils.CommandFlagSpec{Name: "activate", Shorthand: "a", ConfigKey: core.CfgKeyXCommandInstallActivate, Default: "false"},
		)
	})

	Context("command", func() {
		var (
			harness *testutils.CommandHarness
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDownload)

			harness.Cfg.Set(core.CfgKeyXCommandInstallName, "cmdr")
			harness.Cfg.Set(core.CfgKeyXCommandInstallVersion, "1.0.0")
			harness.Cfg.Set(core.CfgKeyXCommandInstallLocation, "")
		})

		AfterEach(func() {
			harness.Finish()
		})

		DescribeTable("should install commands",
			func(activate bool) {
				harness.Cfg.Set(core.CfgKeyXCommandInstallActivate, activate)

				harness.Manager.EXPECT().Define("cmdr", "1.0.0", "")
				if activate {
					harness.Manager.EXPECT().Activate("cmdr", "1.0.0").Return(nil)
				}
				harness.ExpectClose()

				InstallCmd.Run(InstallCmd, []string{})
			},
			Entry("and activates it", true),
			Entry("without activation", false),
		)

		It("should panic when install fails", func() {
			harness.Cfg.Set(core.CfgKeyXCommandInstallActivate, false)
			harness.Manager.EXPECT().Define("cmdr", "1.0.0", "").Return(nil, fmt.Errorf("install failed"))
			harness.ExpectClose()

			Expect(func() { InstallCmd.Run(InstallCmd, []string{}) }).To(Panic())
		})

		It("should change link mode", func() {
			InstallCmd.PreRun(DefineCmd, []string{})

			Expect(harness.Cfg.GetString(core.CfgKeyCmdrLinkMode)).To(Equal("default"))
		})
	})
})
