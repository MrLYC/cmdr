package runner_test

import (
	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/runner"
)

var _ = Describe("CommandDefine", func() {
	var (
		suite   commandTestSuite
		definer define.Runner
	)

	BeforeEach(func() {
		suite.Setup()
		suite.cfg.Set(config.CfgKeyCommandDefineName, suite.name)
		suite.cfg.Set(config.CfgKeyCommandDefineVersion, suite.version)
		suite.cfg.Set(config.CfgKeyCommandDefineLocation, suite.location)

		definer = runner.NewDefineRunner(suite.cfg, suite.helper)
	})

	AfterEach(func() {
		suite.TearDown()
	})

	checkResult := func() {
		suite.WithDB(func(db define.DBClient) {
			var command model.Command
			cnt, err := db.Select(
				q.Eq("Name", suite.name),
				q.Eq("Version", suite.version),
				q.Eq("Location", suite.location),
			).Count(&command)
			Expect(err).To(BeNil())

			Expect(cnt).To(Equal(1))
		})
	}

	It("should define a command", func() {
		Expect(definer.Run(suite.ctx)).To(Succeed())

		checkResult()
	})

	It("should update the command location", func() {
		suite.WithDB(func(db define.DBClient) {
			command := model.Command{
				Name:     suite.name,
				Version:  suite.version,
				Location: "not_exists",
			}

			Expect(db.Save(&command)).To(Succeed())
		})

		Expect(definer.Run(suite.ctx)).To(Succeed())
		checkResult()
	})
})
