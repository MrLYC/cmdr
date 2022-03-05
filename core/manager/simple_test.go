package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

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
		ctrl       *gomock.Controller
		mainMgr    *mock.MockCommandManager
		recoderMgr *mock.MockCommandManager
		mgr        *manager.SimpleManager
		name       = "command"
		version    = "1.0.0"
		location   = "location"
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mainMgr = mock.NewMockCommandManager(ctrl)
		recoderMgr = mock.NewMockCommandManager(ctrl)
		mgr = manager.NewSimpleManager(mainMgr, recoderMgr)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Close", func() {
		It("should close main and recorder manager", func() {
			mainMgr.EXPECT().Close()
			recoderMgr.EXPECT().Close()

			Expect(mgr.Close()).To(Succeed())
		})

		It("should continue close even if fail", func() {
			mainMgr.EXPECT().Close().Return(fmt.Errorf("testing"))
			recoderMgr.EXPECT().Close().Return(fmt.Errorf("testing"))

			Expect(mgr.Close()).NotTo(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			mainMgr.EXPECT().Close().DoAndReturn(func() error {
				ordering = append(ordering, "main")
				return nil
			})
			recoderMgr.EXPECT().Close().DoAndReturn(func() error {
				ordering = append(ordering, "recoder")
				return nil
			})

			Expect(mgr.Close()).To(Succeed())
			Expect(ordering).To(Equal([]string{"main", "recoder"}))
		})
	})

	It("should return provider", func() {
		Expect(mgr.Provider()).To(Equal(core.CommandProviderDefault))
	})

	It("should return query by recorder manager", func() {
		query := mock.NewMockCommandQuery(ctrl)
		recoderMgr.EXPECT().Query().Return(query, nil)

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
			mainMgr.EXPECT().Define(name, version, location).Return(command, nil)
			recoderMgr.EXPECT().Define(name, version, command.GetLocation())

			_, err := mgr.Define(name, version, location)
			Expect(err).To(BeNil())
		})

		It("should call in a specific order", func() {
			var ordering []string

			mainMgr.EXPECT().Define(name, version, location).DoAndReturn(func(name string, version string, location string) (core.Command, error) {
				ordering = append(ordering, "main")
				return command, nil
			})
			recoderMgr.EXPECT().Define(name, version, command.GetLocation()).DoAndReturn(func(name string, version string, location string) (core.Command, error) {
				ordering = append(ordering, "recoder")
				return command, nil
			})

			_, err := mgr.Define(name, version, location)
			Expect(err).To(BeNil())
			Expect(ordering).To(Equal([]string{"main", "recoder"}))
		})

		It("should return when catch error", func() {
			mainMgr.EXPECT().Define(name, version, location).Return(nil, fmt.Errorf("testing"))

			_, err := mgr.Define(name, version, location)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("Undefine", func() {
		It("should call all managers", func() {
			mainMgr.EXPECT().Undefine(name, version)
			recoderMgr.EXPECT().Undefine(name, version)

			Expect(mgr.Undefine(name, version)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			mainMgr.EXPECT().Undefine(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "main")
				return nil
			})
			recoderMgr.EXPECT().Undefine(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "recoder")
				return nil
			})

			Expect(mgr.Undefine(name, version)).To(Succeed())
			Expect(ordering).To(Equal([]string{"recoder", "main"}))
		})

		It("should return when catch error", func() {
			recoderMgr.EXPECT().Undefine(name, version).Return(fmt.Errorf("testing"))

			Expect(mgr.Undefine(name, version)).NotTo(Succeed())
		})
	})

	Context("Activate", func() {
		It("should call all managers", func() {
			mainMgr.EXPECT().Activate(name, version).Return(nil)
			recoderMgr.EXPECT().Activate(name, version).Return(nil)

			Expect(mgr.Activate(name, version)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			mainMgr.EXPECT().Activate(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "main")
				return nil
			})
			recoderMgr.EXPECT().Activate(name, version).DoAndReturn(func(name string, version string) error {
				ordering = append(ordering, "recoder")
				return nil
			})

			Expect(mgr.Activate(name, version)).To(Succeed())
			Expect(ordering).To(Equal([]string{"main", "recoder"}))
		})

		It("should return when catch error", func() {
			mainMgr.EXPECT().Activate(name, version).Return(fmt.Errorf("testing"))

			Expect(mgr.Activate(name, version)).NotTo(Succeed())
		})
	})

	Context("Deactivate", func() {
		It("should call all managers", func() {
			mainMgr.EXPECT().Deactivate(name).Return(nil)
			recoderMgr.EXPECT().Deactivate(name).Return(nil)

			Expect(mgr.Deactivate(name)).To(Succeed())
		})

		It("should call in a specific order", func() {
			var ordering []string

			mainMgr.EXPECT().Deactivate(name).DoAndReturn(func(name string) error {
				ordering = append(ordering, "main")
				return nil
			})
			recoderMgr.EXPECT().Deactivate(name).DoAndReturn(func(name string) error {
				ordering = append(ordering, "recoder")
				return nil
			})

			Expect(mgr.Deactivate(name)).To(Succeed())
			Expect(ordering).To(Equal([]string{"recoder", "main"}))
		})

		It("should return when catch error", func() {
			recoderMgr.EXPECT().Deactivate(name).Return(fmt.Errorf("testing"))

			Expect(mgr.Deactivate(name)).NotTo(Succeed())
		})
	})

	Context("Factory", func() {
		var (
			cfg     core.Configuration
			rootDir string
		)

		BeforeEach(func() {
			cfg = viper.New()

			var err error
			rootDir, err = os.MkdirTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			cfg.Set(core.CfgKeyCmdrDatabasePath, filepath.Join(rootDir, "cmdr.db"))
		})

		AfterEach(func() {
			os.RemoveAll(rootDir)
		})

		It("should create download manager", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderDefault, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.SimpleManager)
			Expect(ok).To(BeTrue())
		})
	})
})
