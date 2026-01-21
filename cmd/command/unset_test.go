package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Unset", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(UnsetCmd, "name", "n", core.CfgKeyXCommandUnsetName, "", true)
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

			cfg.Set(core.CfgKeyXCommandUnsetName, "testing")
		})

		AfterEach(func() {
			ctrl.Finish()
			core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
			core.SetConfiguration(rawCfg)
		})

		It("should unset a command", func() {
			manager.EXPECT().Deactivate("testing").Return(nil)
			manager.EXPECT().Close().Return(nil)

			UnsetCmd.Run(UnsetCmd, []string{})
		})

		It("should not unset a cmdr", func() {
			cfg.Set(core.CfgKeyXCommandUnsetName, "cmdr")
			manager.EXPECT().Close().Return(nil)

			UnsetCmd.Run(UnsetCmd, []string{})
		})
	})
})
