package utils_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Cmdr", func() {
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

	It("should define a activated command", func() {
		manager.EXPECT().Define("test", "1.0.0", "test").Return(nil)
		manager.EXPECT().Activate("test", "1.0.0").Return(nil)

		Expect(utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", true)).To(Succeed())
	})

	It("should define a non-activated command", func() {
		manager.EXPECT().Define("test", "1.0.0", "test").Return(nil)

		Expect(utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", false)).To(Succeed())
	})
})
