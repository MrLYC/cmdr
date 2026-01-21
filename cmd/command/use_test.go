package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Use", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(UseCmd, "name", "n", core.CfgKeyXCommandUseName, "", true)
		testutils.CheckCommandFlag(UseCmd, "version", "v", core.CfgKeyXCommandUseVersion, "", true)
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

			cfg.Set(core.CfgKeyXCommandUseName, "cmdr")
			cfg.Set(core.CfgKeyXCommandUseVersion, "1.0.0")
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should unset a command", func() {
			manager.EXPECT().Activate("cmdr", "1.0.0").Return(nil)
			manager.EXPECT().Close().Return(nil)

			UseCmd.Run(UnsetCmd, []string{})
		})
	})
})
