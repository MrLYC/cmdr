package runner_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asdine/storm/v3/q"
	"github.com/jaswdr/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/model"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

type commandTestSuite struct {
	ctx     context.Context
	cfg     define.Configuration
	helper  *utils.CmdrHelper
	faker   faker.Faker
	command model.Command
}

func (s *commandTestSuite) Setup() {
	tempDir, err := os.MkdirTemp("", "")
	Expect(err).To(BeNil())

	s.ctx = context.Background()
	s.cfg = viper.New()
	s.helper = utils.NewCmdrHelper(tempDir)
	s.faker = faker.New()
	s.command = model.Command{
		Name:     s.faker.Color().ColorName(),
		Version:  s.faker.App().Version(),
		Location: filepath.Join(tempDir, s.faker.Color().ColorName()),
	}

	r := runner.NewMigrateRunner(s.cfg, s.helper)
	Expect(r.Run(s.ctx)).To(Succeed())

	s.MakeCommandBinary()
}

func (s *commandTestSuite) TearDown() {
	Expect(os.RemoveAll(s.helper.GetRootDir())).To(Succeed())
}

func (s *commandTestSuite) Bootstrap() {
	BeforeEach(s.Setup)

	AfterEach(s.TearDown)
}

func (s *commandTestSuite) WithDB(fn func(db define.DBClient)) {
	db, err := core.NewDBClient(s.helper.GetDatabasePath())
	Expect(err).To(BeNil())
	defer db.Close()

	fn(db)
}

func (s *commandTestSuite) MustGetCommandsBy(matchers ...q.Matcher) (commands []*model.Command) {
	s.WithDB(func(db define.DBClient) {
		Expect(db.Select(matchers...).Find(&commands)).To(Succeed())
	})

	return
}

func (s *commandTestSuite) MustGetCommand() *model.Command {
	commands := s.MustGetCommandsBy(
		q.Eq("Name", s.command.Name),
		q.Eq("Version", s.command.Version),
	)

	Expect(len(commands)).To(Equal(1))
	return commands[0]
}

func (s *commandTestSuite) CommandMustNotExists() {
	s.WithDB(func(db define.DBClient) {
		cnt, err := db.Select(q.Eq("Name", s.command.Name), q.Eq("Version", s.command.Version)).Count(&s.command)
		Expect(cnt).To(Equal(0))
		Expect(err).To(BeNil())
	})
}

func (s *commandTestSuite) GetBinaryContent(command *model.Command) []byte {
	return []byte(fmt.Sprintf("#!/bin/sh\necho %s:%s", command.Name, command.Version))
}

func (s *commandTestSuite) MakeCommandBinary() {
	Expect(os.WriteFile(s.command.Location, s.GetBinaryContent(&s.command), 0755)).To(Succeed())
}

func (s *commandTestSuite) RemoveCommandBinary() {
	Expect(os.Remove(s.command.Location)).To(Succeed())
}

func (s *commandTestSuite) UpdateCommandVersion() string {
	version := s.command.Version
	for version == s.command.Version {
		s.command.Version = s.faker.App().Version()
	}

	return s.command.Version
}

func (s *commandTestSuite) UpdateCommandLocation() string {
	location := s.command.Location
	for location == s.command.Location {
		s.command.Location = filepath.Join(s.helper.GetRootDir(), s.faker.Color().ColorName())
	}

	s.MakeCommandBinary()
	return s.command.Location
}

func (s *commandTestSuite) CheckCommandShims(command *model.Command) {
	content, err := os.ReadFile(command.Location)
	Expect(err).To(BeNil())
	Expect(content).To(Equal(s.GetBinaryContent(command)))
}

func (s *commandTestSuite) CheckCommandBin(command *model.Command) {
	shims, err := os.ReadFile(command.Location)
	Expect(err).To(BeNil())

	bin, err := os.ReadFile(s.helper.GetCommandBinPath(s.command.Name))
	Expect(err).To(BeNil())

	Expect(shims).To(Equal(bin))
}

func (s *commandTestSuite) InstallCommand() {
	cfg := viper.New()
	cfg.Set(config.CfgKeyCommandInstallName, s.command.Name)
	cfg.Set(config.CfgKeyCommandInstallVersion, s.command.Version)
	cfg.Set(config.CfgKeyCommandInstallLocation, s.command.Location)

	installer := runner.NewInstallRunner(cfg, s.helper)
	Expect(installer.Run(s.ctx)).To(Succeed())
}

func (s *commandTestSuite) InstallActivatedCommand() {
	cfg := viper.New()
	cfg.Set(config.CfgKeyCommandInstallName, s.command.Name)
	cfg.Set(config.CfgKeyCommandInstallVersion, s.command.Version)
	cfg.Set(config.CfgKeyCommandInstallLocation, s.command.Location)
	cfg.Set(config.CfgKeyCommandInstallActivate, true)

	installer := runner.NewInstallRunner(cfg, s.helper)
	Expect(installer.Run(s.ctx)).To(Succeed())
}
