package runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandUninstall", func() {
	var (
		suite commandTestSuite
	)

	suite.Bootstrap()

	BeforeEach(func() {
		suite.cfg.Set(config.CfgKeyCommandUninstallName, suite.command.Name)
		suite.cfg.Set(config.CfgKeyCommandUninstallVersion, suite.command.Version)
	})

	installCommand := func() {
		installer := runner.NewInstallRunner(suite.cfg, suite.helper)
		Expect(installer.Run(suite.ctx)).To(Succeed())
	}

	Context("Success", func() {
		BeforeEach(func() {
			suite.cfg.Set(config.CfgKeyCommandInstallName, suite.command.Name)
			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.command.Version)
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.command.Location)
		})

		checkResult := func() {
			suite.CommandMustNotExists()
			Expect(suite.helper.GetCommandBinPath(suite.command.Name)).NotTo(BeAnExistingFile())
			Expect(suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version)).NotTo(BeAnExistingFile())
		}

		It("should success to uninstall not activated command", func() {
			installCommand()

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeFalse())

			uninstaller := runner.NewUninstallRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			checkResult()
		})

		It("should success to uninstall even shims not exists", func() {
			installCommand()

			Expect(define.FS.Remove(suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version))).To(Succeed())

			uninstaller := runner.NewUninstallRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			checkResult()
		})

		It("should success to uninstall activated command", func() {
			suite.cfg.Set(config.CfgKeyCommandInstallActivate, true)
			installCommand()

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())

			uninstaller := runner.NewUninstallRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			checkResult()
		})

		It("should success to uninstall even bin not exists", func() {
			suite.cfg.Set(config.CfgKeyCommandInstallActivate, true)
			installCommand()

			Expect(define.FS.Remove(suite.helper.GetCommandBinPath(suite.command.Name))).To(Succeed())

			uninstaller := runner.NewUninstallRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			checkResult()
		})
	})

	Context("Fail", func() {
		It("should fail because command not exists", func() {
			uninstaller := runner.NewUninstallRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).NotTo(Succeed())
		})
	})
})
