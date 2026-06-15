package manager_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/internal/testutils"
	"github.com/mrlyc/cmdr/core/manager"
)

type factoryTestManager struct {
	provider core.CommandProvider
}

func (m factoryTestManager) Close() error { return nil }
func (m factoryTestManager) Provider() core.CommandProvider {
	return m.provider
}
func (m factoryTestManager) Query() (core.CommandQuery, error) { return nil, nil }
func (m factoryTestManager) Define(string, string, string) (core.Command, error) {
	return nil, nil
}
func (m factoryTestManager) Undefine(string, string) error { return nil }
func (m factoryTestManager) Activate(string, string) error { return nil }
func (m factoryTestManager) Deactivate(string) error       { return nil }

var _ = Describe("Factories", func() {
	It("should create binary managers from config", func() {
		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrBinDir, "bin")
		cfg.Set(core.CfgKeyCmdrShimsDir, "shims")
		cfg.Set(core.CfgKeyCmdrLinkMode, "copy")

		mgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(mgr).To(BeAssignableToTypeOf(&manager.BinaryManager{}))

		initializer, err := core.NewInitializer("binary", cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(initializer).To(BeAssignableToTypeOf(&manager.BinaryManager{}))
	})

	It("should create simple, doctor and download managers", func() {
		restoreDatabase := testutils.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(core.Configuration) (core.CommandManager, error) {
			return factoryTestManager{provider: core.CommandProviderDatabase}, nil
		})
		defer restoreDatabase()
		restoreBinary := testutils.RegisterCommandManagerFactory(core.CommandProviderBinary, func(core.Configuration) (core.CommandManager, error) {
			return factoryTestManager{provider: core.CommandProviderBinary}, nil
		})
		defer restoreBinary()

		cfg := viper.New()
		defaultMgr, err := core.NewCommandManager(core.CommandProviderDefault, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(defaultMgr).To(BeAssignableToTypeOf(&manager.SimpleManager{}))

		doctorMgr, err := core.NewCommandManager(core.CommandProviderDoctor, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(doctorMgr).To(BeAssignableToTypeOf(&manager.DoctorManager{}))

		downloadMgr, err := core.NewCommandManager(core.CommandProviderDownload, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(downloadMgr).To(BeAssignableToTypeOf(&manager.DownloadManager{}))
	})

	It("should return factory dependency errors", func() {
		restoreBinary := testutils.RegisterCommandManagerFactory(core.CommandProviderBinary, func(core.Configuration) (core.CommandManager, error) {
			return nil, fmt.Errorf("binary failed")
		})
		defer restoreBinary()
		restoreDatabase := testutils.RestoreCommandManagerFactory(core.CommandProviderDatabase)
		defer restoreDatabase()
		_, err := core.NewCommandManager(core.CommandProviderDoctor, viper.New())
		Expect(err).To(HaveOccurred())

		core.RegisterCommandManagerFactory(core.CommandProviderBinary, func(core.Configuration) (core.CommandManager, error) {
			return factoryTestManager{provider: core.CommandProviderBinary}, nil
		})
		core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(core.Configuration) (core.CommandManager, error) {
			return nil, fmt.Errorf("database failed")
		})
		_, err = core.NewCommandManager(core.CommandProviderDoctor, viper.New())
		Expect(err).To(HaveOccurred())
	})
})
