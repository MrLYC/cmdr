package runner_test

import (
	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandInstall", func() {
	var (
		suite     commandTestSuite
		installer define.Runner
	)

	BeforeEach(func() {
		suite.Setup()
		suite.cfg.Set(config.CfgKeyCommandInstallName, suite.name)
		suite.cfg.Set(config.CfgKeyCommandInstallVersion, suite.version)
		suite.cfg.Set(config.CfgKeyCommandInstallLocation, suite.location)

		Expect(afero.WriteFile(define.FS, suite.location, []byte(`#!/bin/sh\necho $@`), 0755)).To(Succeed())
		installer = runner.NewInstallRunner(suite.cfg, suite.helper)
	})

	AfterEach(func() {
		suite.TearDown()
	})

	It("should install a command", func() {
		Expect(installer.Run(suite.ctx)).To(Succeed())

		suite.WithDB(func(db define.DBClient) {
			var command model.Command
			Expect(db.Select(
				q.Eq("Name", suite.name),
				q.Eq("Version", suite.version),
			).First(&command)).To(BeNil())

			exists, err := afero.Exists(define.FS, command.Location)
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
			Expect(command.Location).To(Equal(suite.helper.GetCommandShimsPath(suite.name, suite.version)))
			Expect(command.Activated).To(BeFalse())
			Expect(command.Managed).To(BeTrue())
		})
	})
})
