package utils_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/internal/testutils"
	"github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Cmdr", func() {
	var (
		ctrl    *gomock.Controller
		manager *mock.MockCommandManager
		query   *mock.MockCommandQuery
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		manager = mock.NewMockCommandManager(ctrl)

		query = mock.NewMockCommandQuery(ctrl)
		manager.EXPECT().Query().Return(query, nil).AnyTimes()

		query.EXPECT().WithName(gomock.Any()).Return(query).AnyTimes()
		query.EXPECT().WithVersion(gomock.Any()).Return(query).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("DefineCmdrCommand", func() {
		It("should define a activated command", func() {
			manager.EXPECT().Define("test", "1.0.0", "test")
			manager.EXPECT().Activate("test", "1.0.0").Return(nil)

			Expect(utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", true)).To(Succeed())
		})

		It("should define a non-activated command", func() {
			manager.EXPECT().Define("test", "1.0.0", "test")

			Expect(utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", false)).To(Succeed())
		})

		It("should return invalid version, define and activate errors", func() {
			_, err := utils.DefineCmdrCommand(manager, "test", "bad", "test", false)
			Expect(err).To(HaveOccurred())

			manager.EXPECT().Define("test", "1.0.0", "test").Return(nil, fmt.Errorf("define failed"))
			_, err = utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", false)
			Expect(err).To(MatchError("define failed"))

			manager.EXPECT().Define("test", "1.0.0", "test")
			manager.EXPECT().Activate("test", "1.0.0").Return(fmt.Errorf("activate failed"))
			_, err = utils.DefineCmdrCommand(manager, "test", "1.0.0", "test", true)
			Expect(err).To(MatchError("activate failed"))
		})
	})

	Context("DefineCmdrCommandNX", func() {
		It("should define a command if not exist", func() {
			query.EXPECT().One().Return(nil, fmt.Errorf("not found"))
			manager.EXPECT().Define("test", "1.0.0", "test")

			Expect(utils.DefineCmdrCommandNX(manager, "test", "1.0.0", "test", false)).To(Succeed())
		})

		It("should not define a command if exist", func() {
			query.EXPECT().One()

			_, err := utils.DefineCmdrCommandNX(manager, "test", "1.0.0", "test", false)
			Expect(errors.Cause(err)).To(Equal(utils.ErrCmdrCommandAlreadyDefined))
		})

		It("should return query errors", func() {
			managerWithQueryError := mock.NewMockCommandManager(ctrl)
			managerWithQueryError.EXPECT().Query().Return(nil, fmt.Errorf("query failed"))

			_, err := utils.DefineCmdrCommandNX(managerWithQueryError, "test", "1.0.0", "test", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("query command"))
		})
	})

	Context("GetCmdrCommand", func() {
		It("should get a command", func() {
			mockCommand := mock.NewMockCommand(ctrl)
			query.EXPECT().One().Return(mockCommand, nil)

			command, err := utils.GetCmdrCommand(manager, "test", "1.0.0")
			Expect(err).To(BeNil())
			Expect(command).To(Equal(mockCommand))
		})

		It("should not get a command", func() {
			query.EXPECT().One().Return(nil, fmt.Errorf("not found"))

			command, err := utils.GetCmdrCommand(manager, "test", "1.0.0")
			Expect(err).NotTo(BeNil())
			Expect(command).To(BeNil())
		})

		It("should return query creation errors", func() {
			managerWithQueryError := mock.NewMockCommandManager(ctrl)
			managerWithQueryError.EXPECT().Query().Return(nil, fmt.Errorf("query failed"))

			command, err := utils.GetCmdrCommand(managerWithQueryError, "test", "1.0.0")
			Expect(err).To(HaveOccurred())
			Expect(command).To(BeNil())
		})
	})

	Context("UpgradeCmdr", func() {
		var (
			ctx            context.Context
			restoreFactory func()
			command        *mock.MockCommand
			url            = "https://example.com"
		)

		BeforeEach(func() {
			ctx = context.Background()
			command = mock.NewMockCommand(ctrl)
			restoreFactory = testutils.RegisterCommandManagerFactory(core.CommandProviderDownload, func(cfg core.Configuration) (core.CommandManager, error) {
				return manager, nil
			})
		})

		AfterEach(func() {
			restoreFactory()
		})

		It("should upgrade a command", func() {
			defined := false
			query.EXPECT().One().DoAndReturn(func() (core.Command, error) {
				if !defined {
					return nil, fmt.Errorf("not found")
				}

				return command, nil
			}).Times(2)
			manager.EXPECT().Define(core.Name, "1.0.0", url).DoAndReturn(func(name, version, url string) (core.Command, error) {
				defined = true
				return nil, nil
			})
			manager.EXPECT().Close().Return(nil)
			command.EXPECT().GetLocation().Return("echo")

			Expect(utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})).To(Succeed())
		})

		It("should not upgrade a command ", func() {
			mockCommand := mock.NewMockCommand(ctrl)
			query.EXPECT().One().Return(mockCommand, nil)

			err := utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})
			Expect(errors.Cause(err)).To(Equal(utils.ErrCmdrCommandAlreadyDefined))
		})

		It("should skip current version upgrades", func() {
			err := utils.UpgradeCmdr(ctx, nil, url, core.Version, []string{})
			Expect(errors.Cause(err)).To(Equal(utils.ErrCmdrAlreadyLatestVersion))
		})

		It("should return manager factory errors", func() {
			restoreFactory()
			restoreFactory = testutils.RegisterCommandManagerFactory(core.CommandProviderDownload, func(cfg core.Configuration) (core.CommandManager, error) {
				return nil, fmt.Errorf("factory failed")
			})

			err := utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("create command manager"))
		})

		It("should return get command errors", func() {
			query.EXPECT().One().Return(nil, fmt.Errorf("not found")).Times(2)
			manager.EXPECT().Define(core.Name, "1.0.0", url).Return(nil, nil)

			err := utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("get command"))
		})

		It("should return close errors", func() {
			defined := false
			query.EXPECT().One().DoAndReturn(func() (core.Command, error) {
				if !defined {
					return nil, fmt.Errorf("not found")
				}
				return command, nil
			}).Times(2)
			manager.EXPECT().Define(core.Name, "1.0.0", url).DoAndReturn(func(name, version, url string) (core.Command, error) {
				defined = true
				return nil, nil
			})
			manager.EXPECT().Close().Return(fmt.Errorf("close failed"))

			err := utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("close command manager"))
		})

		It("should return process errors", func() {
			defined := false
			query.EXPECT().One().DoAndReturn(func() (core.Command, error) {
				if !defined {
					return nil, fmt.Errorf("not found")
				}
				return command, nil
			}).Times(2)
			manager.EXPECT().Define(core.Name, "1.0.0", url).DoAndReturn(func(name, version, url string) (core.Command, error) {
				defined = true
				return nil, nil
			})
			manager.EXPECT().Close().Return(nil)
			command.EXPECT().GetLocation().Return("definitely-not-a-cmdr-test-command")

			err := utils.UpgradeCmdr(ctx, nil, url, "1.0.0", []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("run command"))
		})
	})
})
