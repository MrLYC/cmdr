package command

import (
	. "github.com/onsi/ginkgo"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("List", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(ListCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandListName},
			testutils.CommandFlagSpec{Name: "version", Shorthand: "v", ConfigKey: core.CfgKeyXCommandListVersion},
			testutils.CommandFlagSpec{Name: "location", Shorthand: "l", ConfigKey: core.CfgKeyXCommandListLocation},
			testutils.CommandFlagSpec{Name: "activate", Shorthand: "a", ConfigKey: core.CfgKeyXCommandListActivate, Default: "false"},
		)
	})

	Context("command", func() {
		var (
			harness  *testutils.CommandHarness
			query    *mock.MockCommandQuery
			commands []core.Command
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDefault)
			query = harness.QueryReturning(&commands)
			harness.ExpectClose()
		})

		AfterEach(func() {
			harness.Finish()
		})

		It("should list all commands", func() {
			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by name", func() {
			harness.Cfg.Set(core.CfgKeyXCommandListName, "cmdr")
			query.EXPECT().WithName("cmdr").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by version", func() {
			harness.Cfg.Set(core.CfgKeyXCommandListVersion, "1.0.0")
			query.EXPECT().WithVersion("1.0.0").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by location", func() {
			harness.Cfg.Set(core.CfgKeyXCommandListLocation, "/path/to/cmdr")
			query.EXPECT().WithLocation("/path/to/cmdr").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by activate", func() {
			harness.Cfg.Set(core.CfgKeyXCommandListActivate, true)
			query.EXPECT().WithActivated(true).Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should render selected fields for active and inactive commands", func() {
			commands = []core.Command{
				harness.Command("cmdr", "1.0.0", "/bin/cmdr", true),
				harness.Command("other", "2.0.0", "/bin/other", false),
			}
			harness.Cfg.Set(core.CfgKeyXCommandListFields, []string{"activated", "name", "unknown", "location"})

			ListCmd.Run(ListCmd, []string{})
		})
	})
})
