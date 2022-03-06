package utils_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Command", func() {
	Context("CobraCommandCompleteHelper", func() {
		var (
			ctrl               *gomock.Controller
			mockManager        *mock.MockCommandManager
			mockQuery          *mock.MockCommandQuery
			commandA, commandB *mock.MockCommand
			cobraCommand       *cobra.Command
			helper             *utils.CobraCommandCompleteHelper
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockQuery = mock.NewMockCommandQuery(ctrl)
			mockManager = mock.NewMockCommandManager(ctrl)
			cobraCommand = &cobra.Command{}
			helper = utils.NewCobraCommandCompleteHelper(cobraCommand, core.CommandProviderUnknown)

			mockManager.EXPECT().Query().Return(mockQuery, nil).AnyTimes()
			mockQuery.EXPECT().All().Return([]core.Command{commandA, commandB}, nil).AnyTimes()

			core.RegisterCommandManagerFactory(
				core.CommandProviderUnknown,
				func(cfg core.Configuration) (core.CommandManager, error) {
					return mockManager, nil
				},
			)

			commandA = mock.NewMockCommand(ctrl)
			commandA.EXPECT().GetName().Return("command-a").AnyTimes()
			commandA.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
			commandA.EXPECT().GetLocation().Return("/path/to/command-a").AnyTimes()

			commandB = mock.NewMockCommand(ctrl)
			commandB.EXPECT().GetName().Return("command-b").AnyTimes()
			commandB.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
			commandB.EXPECT().GetLocation().Return("/path/to/command-b").AnyTimes()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should register command when flags not set", func() {
			Expect(helper.RegisterAll()).To(BeNil())
		})

		It("should register all command when flags set", func() {
			flags := cobraCommand.Flags()
			flags.String("name", "", "command name")
			flags.String("version", "", "command version")
			flags.String("location", "", "command location")

			Expect(helper.RegisterAll()).To(BeNil())
		})

		Context("Complete", func() {
			BeforeEach(func() {
				mockManager.EXPECT().Close()
			})

			It("should query command name when it is not empty", func() {
				mockQuery.EXPECT().WithName("x").Return(mockQuery)

				cobraCommand.Flags().String("name", "x", "command name")
				Expect(helper.RegisterNameFunc()).To(BeNil())

				helper.GetNameSlice("")
			})

			It("should query command name by prefix", func() {
				cobraCommand.Flags().String("name", "", "command name")
				Expect(helper.RegisterNameFunc()).To(BeNil())

				Expect(helper.GetNameSlice("command")).To(Equal([]string{"command-a", "command-b"}))
				Expect(helper.GetNameSlice("command-a")).To(Equal([]string{"command-a"}))
			})

			It("should query command version when it is not empty", func() {
				mockQuery.EXPECT().WithVersion("x").Return(mockQuery)

				cobraCommand.Flags().String("version", "x", "command version")
				Expect(helper.RegisterVersionFunc()).To(BeNil())

				helper.GetVersionSlice("")
			})

			It("should query command version by prefix", func() {
				cobraCommand.Flags().String("version", "", "command version")
				Expect(helper.RegisterVersionFunc()).To(BeNil())

				Expect(helper.GetVersionSlice("1.0")).To(Equal([]string{"1.0.0", "1.0.1"}))
				Expect(helper.GetVersionSlice("1.0.1")).To(Equal([]string{"1.0.1"}))
			})

			It("should query command location when it is not empty", func() {
				mockQuery.EXPECT().WithLocation("x").Return(mockQuery)

				cobraCommand.Flags().String("location", "x", "command location")
				Expect(helper.RegisterLocationFunc()).To(BeNil())

				helper.GetLocationSlice("")
			})

			It("should query command location by prefix", func() {
				cobraCommand.Flags().String("location", "", "command location")
				Expect(helper.RegisterLocationFunc()).To(BeNil())

				Expect(helper.GetLocationSlice("/path/to/command")).To(Equal(
					[]string{"/path/to/command-a", "/path/to/command-b"},
				))
				Expect(helper.GetLocationSlice("/path/to/command-a")).To(Equal([]string{"/path/to/command-a"}))
			})
		})

	})
})
