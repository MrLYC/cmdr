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
		testutils.CheckCommandFlag(listCmd, "name", "n", core.CfgKeyXCommandListName, "", false)
		testutils.CheckCommandFlag(listCmd, "version", "v", core.CfgKeyXCommandListVersion, "", false)
		testutils.CheckCommandFlag(listCmd, "location", "l", core.CfgKeyXCommandListLocation, "", false)
		testutils.CheckCommandFlag(listCmd, "activate", "a", core.CfgKeyXCommandListActivate, "false", false)
	})

	Context("command", func() {
		var (
			ctrl    *gomock.Controller
			rawCfg  core.Configuration
			cfg     core.Configuration
			manager *mock.MockCommandManager
			query   *mock.MockCommandQuery
			factory func(cfg core.Configuration) (core.CommandManager, error)
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
			query.EXPECT().All().Return(nil, nil)

			manager.EXPECT().Query().Return(query, nil)
			manager.EXPECT().Close().Return(nil)
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should list all commands", func() {
			listCmd.Run(listCmd, []string{})
		})

		It("should filter by name", func() {
			cfg.Set(core.CfgKeyXCommandListName, "cmdr")
			query.EXPECT().WithName("cmdr").Return(nil)

			listCmd.Run(listCmd, []string{})
		})

		It("should filter by version", func() {
			cfg.Set(core.CfgKeyXCommandListVersion, "1.0.0")
			query.EXPECT().WithVersion("1.0.0").Return(nil)

			listCmd.Run(listCmd, []string{})
		})

		It("should filter by location", func() {
			cfg.Set(core.CfgKeyXCommandListLocation, "/path/to/cmdr")
			query.EXPECT().WithLocation("/path/to/cmdr").Return(nil)

			listCmd.Run(listCmd, []string{})
		})

		It("should filter by activate", func() {
			cfg.Set(core.CfgKeyXCommandListActivate, true)
			query.EXPECT().WithActivated(true).Return(nil)

			listCmd.Run(listCmd, []string{})
		})
	})
})
