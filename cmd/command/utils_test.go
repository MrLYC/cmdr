package command

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Utils", func() {
	var (
		ctrl    *gomock.Controller
		manager *mock.MockCommandManager
		factory func(cfg core.Configuration) (core.CommandManager, error)
	)

	BeforeEach(func() {
		factory = core.GetCommandManagerFactory(core.CommandProviderDefault)

		ctrl = gomock.NewController(GinkgoT())
		manager = mock.NewMockCommandManager(ctrl)
		core.RegisterCommandManagerFactory(core.CommandProviderDefault, func(cfg core.Configuration) (core.CommandManager, error) {
			return manager, nil
		})

	})

	AfterEach(func() {
		ctrl.Finish()
		core.RegisterCommandManagerFactory(core.CommandProviderDefault, factory)
	})

	It("should init manager", func() {
		var cmd cobra.Command

		manager.EXPECT().Close().Return(nil)

		fn := runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
			Expect(manager).NotTo(BeNil())

			return nil
		})

		fn(&cmd, []string{})
	})
})
