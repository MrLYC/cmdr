package command

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Utils", func() {
	var (
		harness *testutils.CommandHarness
	)

	BeforeEach(func() {
		harness = testutils.NewCommandHarness(core.CommandProviderDefault)
	})

	AfterEach(func() {
		harness.Finish()
	})

	It("should init manager", func() {
		var cmd cobra.Command

		harness.ExpectClose()

		fn := runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
			Expect(manager).NotTo(BeNil())

			return nil
		})

		fn(&cmd, []string{})
	})

	It("should query commands", func() {
		query := mock.NewMockCommandQuery(harness.Ctrl)
		query.EXPECT().WithActivated(true).Return(query)
		query.EXPECT().WithName("name").Return(query)
		query.EXPECT().WithVersion("version").Return(query)
		query.EXPECT().WithLocation("location").Return(query)
		query.EXPECT().All().Return(nil, nil)

		harness.Manager.EXPECT().Query().Return(query, nil)

		_, err := queryCommands(harness.Manager, true, "name", "version", "location")
		Expect(err).To(BeNil())
	})
})
