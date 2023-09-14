package initializer_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/initializer"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Command", func() {
	var (
		ctrl             *gomock.Controller
		query            *mock.MockCommandQuery
		manager          *mock.MockCommandManager
		legacyCommand    *mock.MockCommand
		activatedCommand *mock.MockCommand
		name             = "cmdr"
		version          = "100.0.0"
		location         = "github.com/mrlyc/cmdr"
		updater          *initializer.CmdrUpdater
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		legacyCommand = mock.NewMockCommand(ctrl)
		legacyCommand.EXPECT().GetName().Return(name).AnyTimes()
		legacyCommand.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
		legacyCommand.EXPECT().GetActivated().Return(false).AnyTimes()

		activatedCommand = mock.NewMockCommand(ctrl)
		activatedCommand.EXPECT().GetName().Return(name).AnyTimes()
		activatedCommand.EXPECT().GetVersion().Return("10.0.0").AnyTimes()
		activatedCommand.EXPECT().GetActivated().Return(true).AnyTimes()

		query = mock.NewMockCommandQuery(ctrl)
		query.EXPECT().WithName(name).Return(query).AnyTimes()
		query.EXPECT().WithActivated(true).Return(query).AnyTimes()
		query.EXPECT().All().Return([]core.Command{legacyCommand, activatedCommand}, nil).AnyTimes()
		query.EXPECT().One().Return(activatedCommand, nil).AnyTimes()
		query.EXPECT().Count().Return(2, nil).AnyTimes()

		manager = mock.NewMockCommandManager(ctrl)
		manager.EXPECT().Query().Return(query, nil).AnyTimes()

		updater = initializer.NewCmdrUpdater(
			manager, name, version, location,
		)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should update cmdr", func() {
		manager.EXPECT().Define(name, version, location)
		manager.EXPECT().Activate(name, version)
		manager.EXPECT().Undefine(name, legacyCommand.GetVersion())

		Expect(updater.Init(false)).To(Succeed())
	})

	It("should upgrade cmdr", func() {
		manager.EXPECT().Activate(name, version)
		manager.EXPECT().Undefine(name, legacyCommand.GetVersion())

		Expect(updater.Init(true)).To(Succeed())
	})
})
