package cmd

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Doctor", func() {
	It("should run doctor command with options", func() {
		previousDryRun := doctorDryRun
		previousNoBackup := doctorNoBackup
		defer func() {
			doctorDryRun = previousDryRun
			doctorNoBackup = previousNoBackup
		}()

		root, err := os.MkdirTemp("", "cmdr-doctor")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(root)

		cfg := core.NewConfiguration()
		cfg.Set(core.CfgKeyCmdrRootDir, filepath.Join(root, "missing"))
		restoreConfig := testutils.SwapConfiguration(cfg)
		defer restoreConfig()
		restoreFactory := testutils.RegisterCommandManagerFactory(core.CommandProviderDoctor, func(cfg core.Configuration) (core.CommandManager, error) {
			return &cleanTestManager{}, nil
		})
		defer restoreFactory()

		doctorDryRun = true
		doctorNoBackup = true
		Expect(func() { doctorCmd.Run(doctorCmd, nil) }).NotTo(Panic())
	})
})
