package runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandUse", func() {
	var (
		suite commandTestSuite
	)

	suite.Bootstrap()

	BeforeEach(func() {
		suite.cfg.Set(config.CfgKeyCommandUseName, suite.command.Name)
		suite.cfg.Set(config.CfgKeyCommandUseVersion, suite.command.Version)
	})

	Context("Success", func() {
		It("should success to use a command", func() {
			suite.InstallCommand()

			use := runner.NewUseRunner(suite.cfg, suite.helper)
			Expect(use.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())
			Expect(suite.helper.GetCommandBinPath(command.Name)).To(BeAnExistingFile())
		})

		It("should success to use an activated command", func() {
			suite.InstallActivatedCommand()

			use := runner.NewUseRunner(suite.cfg, suite.helper)
			Expect(use.Run(suite.ctx)).To(Succeed())

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())
			Expect(suite.helper.GetCommandBinPath(command.Name)).To(BeAnExistingFile())
		})
	})
})
