package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("List", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(ListCmd, "name", "n", core.CfgKeyXCommandListName, "", false)
		testutils.CheckCommandFlag(ListCmd, "version", "v", core.CfgKeyXCommandListVersion, "", false)
		testutils.CheckCommandFlag(ListCmd, "location", "l", core.CfgKeyXCommandListLocation, "", false)
		testutils.CheckCommandFlag(ListCmd, "activate", "a", core.CfgKeyXCommandListActivate, "false", false)
	})

	Context("command", func() {
		var (
			ctrl     *gomock.Controller
			rawCfg   core.Configuration
			cfg      core.Configuration
			manager  *mock.MockCommandManager
			query    *mock.MockCommandQuery
			factory  func(cfg core.Configuration) (core.CommandManager, error)
			commands []core.Command
		)

		BeforeEach(func() {
			factory = core.GetCommandManagerFactory(core.CommandProviderDefault)
			rawCfg = core.GetConfiguration()

			ctrl = gomock.NewController(GinkgoT())
			manager = mock.NewMockCommandManager(ctrl)
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, func(cfg core.Configuration) (core.CommandManager, error) {
				return manager, nil
			})

			cfg = viper.New()
			core.SetConfiguration(cfg)

			query = mock.NewMockCommandQuery(ctrl)
			query.EXPECT().All().DoAndReturn(func() ([]core.Command, error) {
				return commands, nil
			}).AnyTimes()

			manager.EXPECT().Query().Return(query, nil)
			manager.EXPECT().Close().Return(nil)
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should list all commands", func() {
			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by name", func() {
			cfg.Set(core.CfgKeyXCommandListName, "cmdr")
			query.EXPECT().WithName("cmdr").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by version", func() {
			cfg.Set(core.CfgKeyXCommandListVersion, "1.0.0")
			query.EXPECT().WithVersion("1.0.0").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by location", func() {
			cfg.Set(core.CfgKeyXCommandListLocation, "/path/to/cmdr")
			query.EXPECT().WithLocation("/path/to/cmdr").Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should filter by activate", func() {
			cfg.Set(core.CfgKeyXCommandListActivate, true)
			query.EXPECT().WithActivated(true).Return(nil)

			ListCmd.Run(ListCmd, []string{})
		})

		It("should render selected fields for active and inactive commands", func() {
			active := mock.NewMockCommand(ctrl)
			active.EXPECT().GetActivated().Return(true).AnyTimes()
			active.EXPECT().GetName().Return("cmdr").AnyTimes()
			active.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
			active.EXPECT().GetLocation().Return("/bin/cmdr").AnyTimes()
			inactive := mock.NewMockCommand(ctrl)
			inactive.EXPECT().GetActivated().Return(false).AnyTimes()
			inactive.EXPECT().GetName().Return("other").AnyTimes()
			inactive.EXPECT().GetVersion().Return("2.0.0").AnyTimes()
			inactive.EXPECT().GetLocation().Return("/bin/other").AnyTimes()
			commands = []core.Command{active, inactive}
			cfg.Set(core.CfgKeyXCommandListFields, []string{"activated", "name", "unknown", "location"})

			ListCmd.Run(ListCmd, []string{})
		})
	})
})
