package runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandUnset", func() {
	var (
		suite commandTestSuite
	)

	suite.Bootstrap()

	BeforeEach(func() {
		suite.cfg.Set(config.CfgKeyCommandUnsetName, suite.command.Name)
	})

	Context("Success", func() {
		It("should success to uninstall not activated command", func() {
			suite.InstallCommand()

			uninstaller := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeFalse())
		})
		It("should success to uninstall activated command", func() {
			suite.InstallActivatedCommand()

			uninstaller := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeFalse())
		})

		It("should success to uninstall even shims not exists", func() {
			suite.InstallActivatedCommand()

			Expect(define.FS.Remove(suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version))).To(Succeed())

			uninstaller := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())
		})

		It("should success to uninstall even shims not exists", func() {
			suite.InstallActivatedCommand()

			Expect(define.FS.Remove(suite.helper.GetCommandBinPath(suite.command.Name))).To(Succeed())

			uninstaller := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(uninstaller.Run(suite.ctx)).To(Succeed())
		})
	})
})
