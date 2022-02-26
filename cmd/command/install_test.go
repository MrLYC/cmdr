package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Install", func() {
	It("should check flags", func() {
		checkCommandFlag(installCmd, "name", "n", core.CfgKeyXCommandInstallName, "", true)
		checkCommandFlag(installCmd, "version", "v", core.CfgKeyXCommandInstallVersion, "", true)
		checkCommandFlag(installCmd, "location", "l", core.CfgKeyXCommandInstallLocation, "", true)
		checkCommandFlag(installCmd, "activate", "a", core.CfgKeyXCommandInstallActivate, "false", false)
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

			cfg.Set(core.CfgKeyXCommandInstallName, "cmdr")
			cfg.Set(core.CfgKeyXCommandInstallVersion, "1.0.0")
			cfg.Set(core.CfgKeyXCommandInstallLocation, "")
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should install a activated command", func() {
			cfg.Set(core.CfgKeyXCommandInstallActivate, true)

			manager.EXPECT().Define("cmdr", "1.0.0", "").Return(nil)
			manager.EXPECT().Activate("cmdr", "1.0.0").Return(nil)
			manager.EXPECT().Close().Return(nil)

			installCmd.Run(installCmd, []string{})
		})

		It("should install a non-activated command", func() {
			cfg.Set(core.CfgKeyXCommandInstallActivate, false)

			manager.EXPECT().Define("cmdr", "1.0.0", "").Return(nil)
			manager.EXPECT().Close().Return(nil)

			installCmd.Run(installCmd, []string{})
		})

		It("should change link mode", func() {
			installCmd.PreRun(defineCmd, []string{})

			Expect(cfg.GetString(core.CfgKeyCmdrLinkMode)).To(Equal("default"))
		})
	})
})
