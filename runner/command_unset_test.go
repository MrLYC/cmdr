package runner_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
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
		It("should success to unset not activated command", func() {
			suite.InstallCommand()

			unsetter := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(unsetter.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeFalse())
			Expect(suite.helper.GetCommandBinPath(command.Name)).NotTo(BeAnExistingFile())
		})

		It("should success to unset activated command", func() {
			suite.InstallActivatedCommand()

			unsetter := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(unsetter.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeFalse())
			Expect(suite.helper.GetCommandBinPath(command.Name)).NotTo(BeAnExistingFile())
		})

		It("should success to unset even shims not exists", func() {
			suite.InstallActivatedCommand()

			Expect(os.Remove(suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version))).To(Succeed())

			unsetter := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(unsetter.Run(suite.ctx)).To(Succeed())
		})

		It("should success to unset even shims not exists", func() {
			suite.InstallActivatedCommand()

			Expect(os.Remove(suite.helper.GetCommandBinPath(suite.command.Name))).To(Succeed())

			unsetter := runner.NewUnsetRunner(suite.cfg, suite.helper)
			Expect(unsetter.Run(suite.ctx)).To(Succeed())
		})
	})
})
