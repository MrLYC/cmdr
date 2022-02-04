package runner_test

import (
	"path/filepath"

	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandInstall", func() {
	var (
		suite commandTestSuite
	)

	suite.Bootstrap()

	runInstaller := func() {
		installer := runner.NewInstallRunner(suite.cfg, suite.helper)
		Expect(installer.Run(suite.ctx)).To(Succeed())
	}

	Context("Default config", func() {
		BeforeEach(func() {
			suite.cfg.Set(config.CfgKeyCommandInstallName, suite.command.Name)
			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.command.Version)
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.command.Location)
		})

		It("should install a command", func() {
			runInstaller()

			command := suite.MustGetCommand()
			Expect(command.Location).To(Equal(suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version)))
			Expect(command.Activated).To(BeFalse())
			Expect(command.Managed).To(BeTrue())

			suite.CheckCommandShims(command)
		})

		It("should install a command with different version", func() {
			runInstaller()

			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.UpdateCommandVersion())
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.UpdateCommandLocation())
			runInstaller()

			commands := suite.MustGetCommandsBy(q.Eq("Name", suite.command.Name))
			Expect(commands).To(HaveLen(2))

			suite.CheckCommandShims(commands[0])
			suite.CheckCommandShims(commands[1])
		})

		It("should update the command location", func() {
			runInstaller()

			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.UpdateCommandLocation())

			runInstaller()
			command := suite.MustGetCommand()
			suite.CheckCommandShims(command)
		})

		It("should overwrite exists shims", func() {
			shimsPath := suite.helper.GetCommandShimsPath(suite.command.Name, suite.command.Version)
			Expect(define.FS.MkdirAll(filepath.Dir(shimsPath), 0755)).To(Succeed())
			Expect(afero.WriteFile(define.FS, shimsPath, []byte{}, 0755)).To(Succeed())

			runInstaller()
			suite.CheckCommandShims(&suite.command)
		})

		It("should fail because binary not exists", func() {
			suite.RemoveCommandBinary()
			installer := runner.NewInstallRunner(suite.cfg, suite.helper)
			Expect(installer.Run(suite.ctx)).NotTo(Succeed())

			suite.CommandMustNotExists()
		})

		It("should fail because location not exists", func() {
			Expect(define.FS.Remove(suite.command.Location)).To(Succeed())
			installer := runner.NewInstallRunner(suite.cfg, suite.helper)
			Expect(installer.Run(suite.ctx)).NotTo(Succeed())

			suite.CommandMustNotExists()
		})
	})

	Context("Activate mode", func() {
		BeforeEach(func() {
			suite.cfg.Set(config.CfgKeyCommandInstallName, suite.command.Name)
			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.command.Version)
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.command.Location)
			suite.cfg.Set(config.CfgKeyCommandInstallActivate, true)
		})

		It("should activate a command", func() {
			runInstaller()

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())
			suite.CheckCommandBin(command)
		})

		It("should activate a command with different version", func() {
			runInstaller()

			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.UpdateCommandVersion())
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.UpdateCommandLocation())
			runInstaller()

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())
			suite.CheckCommandBin(command)
		})

		It("should activate a command even the bin does not exist", func() {
			runInstaller()
			Expect(define.FS.Remove(suite.helper.GetCommandBinPath(suite.command.Name))).To(Succeed())

			suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.UpdateCommandVersion())
			suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.UpdateCommandLocation())
			runInstaller()

			command := suite.MustGetCommand()
			Expect(command.Activated).To(BeTrue())
			suite.CheckCommandBin(command)
		})
	})
})
