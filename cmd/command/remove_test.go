package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Remove", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(RemoveCmd, "name", "n", core.CfgKeyXCommandRemoveName, "", true)
		testutils.CheckCommandFlag(RemoveCmd, "version", "v", core.CfgKeyXCommandRemoveVersion, "", true)
	})

	Context("command", func() {
		var (
			ctrl    *gomock.Controller
			rawCfg  core.Configuration
			cfg     core.Configuration
			manager *mock.MockCommandManager
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
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should undefine a command", func() {
			cfg.Set(core.CfgKeyXCommandRemoveName, "cmdr")
			cfg.Set(core.CfgKeyXCommandRemoveVersion, "1.0.0")

			manager.EXPECT().Undefine("cmdr", "1.0.0").Return(nil)
			manager.EXPECT().Close().Return(nil)

			RemoveCmd.Run(RemoveCmd, []string{})
		})

		It("should not undefine a activated command", func() {
			cfg.Set(core.CfgKeyXCommandRemoveName, "cmdr")
			cfg.Set(core.CfgKeyXCommandRemoveVersion, "1.0.0")

			manager.EXPECT().Undefine("cmdr", "1.0.0").Return(core.ErrCommandAlreadyActivated)
			manager.EXPECT().Close().Return(nil)

			RemoveCmd.Run(RemoveCmd, []string{})
		})
	})
})
