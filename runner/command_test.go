package runner_test

import (
	"context"
	"path/filepath"

	"github.com/jaswdr/faker"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

type commandTestSuite struct {
	ctx      context.Context
	cfg      define.Configuration
	helper   *utils.CmdrHelper
	faker    faker.Faker
	name     string
	version  string
	location string
}

func (s *commandTestSuite) Setup() {
	tempDir, err := afero.TempDir(define.FS, "", "")
	Expect(err).To(BeNil())

	s.ctx = context.Background()
	s.cfg = viper.New()
	s.helper = utils.NewCmdrHelper(tempDir)
	s.faker = faker.New()
	s.name = s.faker.Color().ColorName()
	s.version = s.faker.App().Version()
	s.location = filepath.Join(tempDir, s.faker.File().FilenameWithExtension())

	r := runner.NewMigrateRunner(s.cfg, s.helper)
	Expect(r.Run(s.ctx)).To(Succeed())
}

func (s *commandTestSuite) TearDown() {
	Expect(define.FS.RemoveAll(s.helper.GetRootDir())).To(Succeed())
}

func (s *commandTestSuite) WithDB(fn func(db define.DBClient)) {
	db, err := operator.NewDBClient(s.helper.GetDatabasePath())
	Expect(err).To(BeNil())
	defer db.Close()

	fn(db)
}
