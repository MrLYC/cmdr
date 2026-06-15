package command

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Remove", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlags(RemoveCmd,
			testutils.CommandFlagSpec{Name: "name", Shorthand: "n", ConfigKey: core.CfgKeyXCommandRemoveName, IsRequired: true},
			testutils.CommandFlagSpec{Name: "version", Shorthand: "v", ConfigKey: core.CfgKeyXCommandRemoveVersion, IsRequired: true},
		)
	})

	Context("command", func() {
		var (
			harness *testutils.CommandHarness
		)

		BeforeEach(func() {
			harness = testutils.NewCommandHarness(core.CommandProviderDefault)

			harness.Cfg.Set(core.CfgKeyXCommandRemoveName, "cmdr")
			harness.Cfg.Set(core.CfgKeyXCommandRemoveVersion, "1.0.0")
		})

		AfterEach(func() {
			harness.Finish()
		})

		DescribeTable("should handle expected undefine results",
			func(err error) {
				harness.Manager.EXPECT().Undefine("cmdr", "1.0.0").Return(err)
				harness.ExpectClose()

				RemoveCmd.Run(RemoveCmd, []string{})
			},
			Entry("successfully", nil),
			Entry("when command is activated", core.ErrCommandAlreadyActivated),
		)

		It("should panic when undefine fails", func() {
			harness.Manager.EXPECT().Undefine("cmdr", "1.0.0").Return(fmt.Errorf("remove failed"))
			harness.ExpectClose()

			Expect(func() { RemoveCmd.Run(RemoveCmd, []string{}) }).To(Panic())
		})
	})
})
