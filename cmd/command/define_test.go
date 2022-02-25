package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Define", func() {
	It("should check flags", func() {
		checkCommandFlag(defineCmd, "name", "n", core.CfgKeyXCommandDefineName, "", true)
		checkCommandFlag(defineCmd, "version", "v", core.CfgKeyXCommandDefineVersion, "", true)
		checkCommandFlag(defineCmd, "location", "l", core.CfgKeyXCommandDefineLocation, "", true)
		checkCommandFlag(defineCmd, "activate", "a", core.CfgKeyXCommandDefineActivate, "false", false)
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

		It("should define a activated command", func() {
			cfg.Set(core.CfgKeyXCommandDefineName, "test")
			cfg.Set(core.CfgKeyXCommandDefineVersion, "1.0.0")
			cfg.Set(core.CfgKeyXCommandDefineLocation, "")
			cfg.Set(core.CfgKeyXCommandDefineActivate, true)

			manager.EXPECT().Define("test", "1.0.0", "").Return(nil)
			manager.EXPECT().Activate("test", "1.0.0").Return(nil)
			manager.EXPECT().Close().Return(nil)

			defineCmd.Run(defineCmd, []string{})
		})

		It("should define a non-activated command", func() {
			cfg.Set(core.CfgKeyXCommandDefineName, "test")
			cfg.Set(core.CfgKeyXCommandDefineVersion, "1.0.0")
			cfg.Set(core.CfgKeyXCommandDefineLocation, "")
			cfg.Set(core.CfgKeyXCommandDefineActivate, false)

			manager.EXPECT().Define("test", "1.0.0", "").Return(nil)
			manager.EXPECT().Close().Return(nil)

			defineCmd.Run(defineCmd, []string{})
		})
	})
})
