package runner_test

import (
	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandDefine", func() {
	var (
		suite commandTestSuite
	)

	suite.Bootstrap()

	runDefiner := func() {
		definer := runner.NewDefineRunner(suite.cfg, suite.helper)
		Expect(definer.Run(suite.ctx)).To(Succeed())
	}

	Context("Default config", func() {
		BeforeEach(func() {
			suite.cfg.Set(config.CfgKeyCommandDefineName, suite.command.Name)
			suite.cfg.Set(config.CfgKeyCommandDefineVersion, suite.command.Version)
			suite.cfg.Set(config.CfgKeyCommandDefineLocation, suite.command.Location)
		})

		It("should define a command", func() {
			runDefiner()

			command := suite.MustGetCommand()
			Expect(command.Location).To(Equal(suite.command.Location))
			Expect(command.Managed).To(BeFalse())
		})

		It("should define a command with different version", func() {
			definer := runner.NewDefineRunner(suite.cfg, suite.helper)
			Expect(definer.Run(suite.ctx)).To(Succeed())

			suite.cfg.Set(config.CfgKeyCommandDefineVersion, suite.UpdateCommandVersion())
			definer = runner.NewDefineRunner(suite.cfg, suite.helper)
			Expect(definer.Run(suite.ctx)).To(Succeed())

			commands := suite.MustGetCommandsBy(q.Eq("Name", suite.command.Name))
			Expect(commands).To(HaveLen(2))
		})

		It("should update the command location", func() {
			runDefiner()
			suite.cfg.Set(config.CfgKeyCommandDefineLocation, suite.UpdateCommandLocation())

			runDefiner()
			command := suite.MustGetCommand()
			Expect(command.Location).To(Equal(suite.command.Location))
		})

		It("should raise error", func() {
			Expect(define.FS.Remove(suite.command.Location)).To(Succeed())
			definer := runner.NewDefineRunner(suite.cfg, suite.helper)
			Expect(definer.Run(suite.ctx)).NotTo(BeNil())

			suite.CommandMustNotExists()
		})

	})
})
