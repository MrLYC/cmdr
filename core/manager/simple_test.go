package manager_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Simple", func() {
	var (
		ctrl        *gomock.Controller
		binaryMgr   *mock.MockCommandManager
		databaseMgr *mock.MockCommandManager
		mgr         *manager.SimpleManager
		name        = "command"
		version     = "1.0.0"
		location    = "location"
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		binaryMgr = mock.NewMockCommandManager(ctrl)
		databaseMgr = mock.NewMockCommandManager(ctrl)
		mgr = manager.NewSimpleManager(databaseMgr, []core.CommandManager{binaryMgr})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Close", func() {
		It("should close main and recorder manager", func() {
			binaryMgr.EXPECT().Close()
			databaseMgr.EXPECT().Close()

			Expect(mgr.Close()).To(Succeed())
		})

		It("should continue close even if fail", func() {
			binaryMgr.EXPECT().Close().Return(fmt.Errorf("testing"))
			databaseMgr.EXPECT().Close().Return(fmt.Errorf("testing"))

			Expect(mgr.Close()).NotTo(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			binaryMgr.EXPECT().Close().DoAndReturn(func() error {
				ordering = append(ordering, "binary")
				return nil
			})
			databaseMgr.EXPECT().Close().DoAndReturn(func() error {
				ordering = append(ordering, "database")
				return nil
			})

			Expect(mgr.Close()).To(Succeed())
			Expect(ordering).To(Equal([]string{"database", "binary"}))
		})
	})

	It("should return provider", func() {
		Expect(mgr.Provider()).To(Equal(core.CommandProviderDefault))
	})

	It("should return query by recorder manager", func() {
		query := mock.NewMockCommandQuery(ctrl)
		databaseMgr.EXPECT().Query().Return(query, nil)

		result, err := mgr.Query()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(query))
	})

	Context("Define", func() {
		var (
			command *mock.MockCommand
		)

		BeforeEach(func() {
			command = mock.NewMockCommand(ctrl)
			command.EXPECT().GetLocation().Return("shims_path").AnyTimes()
		})

		It("should call all managers", func() {
			binaryMgr.EXPECT().Define(name, version, location).Return(command, nil)
			databaseMgr.EXPECT().Define(name, version, location)

			_, err := mgr.Define(name, version, location)
			Expect(err).To(BeNil())
		})

		It("should call in a specific order", func() {
			var ordering []string

			binaryMgr.EXPECT().Define(name, version, location).DoAndReturn(func(name string, version string, location string) (core.Command, error) {
				ordering = append(ordering, "binary")
				return command, nil
			})
			databaseMgr.EXPECT().Define(name, version, location).DoAndReturn(func(name string, version string, location string) (core.Command, error) {
				ordering = append(ordering, "database")
				return command, nil
			})

			_, err := mgr.Define(name, version, location)
			Expect(err).To(BeNil())
			Expect(ordering).To(Equal([]string{"database", "binary"}))
		})

		It("should return when catch error", func() {
			databaseMgr.EXPECT().Define(name, version, location).Return(nil, fmt.Errorf("testing"))

			_, err := mgr.Define(name, version, location)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("Undefine", func() {
		It("should call all managers", func() {
			binaryMgr.EXPECT().Undefine(name, version)
			databaseMgr.EXPECT().Undefine(name, version)

			Expect(mgr.Undefine(name, version)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			binaryMgr.EXPECT().Undefine(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "main")
				return nil
			})
			databaseMgr.EXPECT().Undefine(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "recorder")
				return nil
			})

			Expect(mgr.Undefine(name, version)).To(Succeed())
			Expect(ordering).To(Equal([]string{"recorder", "main"}))
		})

		It("should return when catch error", func() {
			databaseMgr.EXPECT().Undefine(name, version).Return(fmt.Errorf("testing"))

			Expect(mgr.Undefine(name, version)).NotTo(Succeed())
		})
	})

	Context("Activate", func() {
		It("should call all managers", func() {
			binaryMgr.EXPECT().Activate(name, version).Return(nil)
			databaseMgr.EXPECT().Activate(name, version).Return(nil)

			Expect(mgr.Activate(name, version)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			binaryMgr.EXPECT().Activate(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "binary")
				return nil
			})
			databaseMgr.EXPECT().Activate(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "database")
				return nil
			})

			Expect(mgr.Activate(name, version)).To(Succeed())
			Expect(ordering).To(Equal([]string{"database", "binary"}))
		})

		It("should return when catch error", func() {
			databaseMgr.EXPECT().Activate(name, version).Return(fmt.Errorf("testing"))

			Expect(mgr.Activate(name, version)).NotTo(Succeed())
		})
	})

	Context("Deactivate", func() {
		It("should call all managers", func() {
			binaryMgr.EXPECT().Deactivate(name).Return(nil)
			databaseMgr.EXPECT().Deactivate(name).Return(nil)

			Expect(mgr.Deactivate(name)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			binaryMgr.EXPECT().Deactivate(name).DoAndReturn(func(name string) error {
				ordering = append(ordering, "binary")
				return nil
			})
			databaseMgr.EXPECT().Deactivate(name).DoAndReturn(func(name string) error {
				ordering = append(ordering, "database")
				return nil
			})

			Expect(mgr.Deactivate(name)).To(Succeed())
			Expect(ordering).To(Equal([]string{"database", "binary"}))
		})

		It("should return when catch error", func() {
			databaseMgr.EXPECT().Deactivate(name).Return(fmt.Errorf("testing"))

			Expect(mgr.Deactivate(name)).NotTo(Succeed())
		})
	})

	Context("Factory", func() {
		var (
			cfg       core.Configuration
			db        *mock.MockDatabase
			dbFactory func() (core.Database, error)
		)

		BeforeEach(func() {
			cfg = viper.New()
			db = mock.NewMockDatabase(ctrl)
			core.SetDatabaseFactory(func() (core.Database, error) {
				return db, nil
			})
			dbFactory = core.GetDatabaseFactory()
		})

		AfterEach(func() {
			core.SetDatabaseFactory(dbFactory)
		})

		It("should create download manager", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderDefault, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.SimpleManager)
			Expect(ok).To(BeTrue())
		})
	})
})
