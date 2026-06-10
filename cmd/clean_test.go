package cmd

import (
	"os"
	"path/filepath"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

type cleanTestCommand struct {
	name      string
	version   string
	activated bool
	location  string
}

func (c cleanTestCommand) GetName() string {
	return c.name
}

func (c cleanTestCommand) GetVersion() string {
	return c.version
}

func (c cleanTestCommand) GetActivated() bool {
	return c.activated
}

func (c cleanTestCommand) GetLocation() string {
	return c.location
}

type cleanTestQuery struct {
	commands []core.Command
	err      error
}

func (q cleanTestQuery) WithName(string) core.CommandQuery {
	return q
}

func (q cleanTestQuery) WithVersion(string) core.CommandQuery {
	return q
}

func (q cleanTestQuery) WithActivated(bool) core.CommandQuery {
	return q
}

func (q cleanTestQuery) WithLocation(string) core.CommandQuery {
	return q
}

func (q cleanTestQuery) All() ([]core.Command, error) {
	return q.commands, q.err
}

func (q cleanTestQuery) One() (core.Command, error) {
	return nil, nil
}

func (q cleanTestQuery) Count() (int, error) {
	return len(q.commands), nil
}

type cleanTestManager struct {
	commands  []core.Command
	queryErr  error
	allErr    error
	undefErr  error
	undefined []string
}

func (m *cleanTestManager) Close() error {
	return nil
}

func (m *cleanTestManager) Provider() core.CommandProvider {
	return core.CommandProviderDefault
}

func (m *cleanTestManager) Query() (core.CommandQuery, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return cleanTestQuery{commands: m.commands, err: m.allErr}, nil
}

func (m *cleanTestManager) Define(string, string, string) (core.Command, error) {
	return nil, nil
}

func (m *cleanTestManager) Undefine(name string, version string) error {
	m.undefined = append(m.undefined, name+":"+version)
	return m.undefErr
}

func (m *cleanTestManager) Activate(string, string) error {
	return nil
}

func (m *cleanTestManager) Deactivate(string) error {
	return nil
}

var _ = Describe("Clean", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(cleanCmd, "age", "", core.CfgKeyXCleanAgeDays, "100", false)
		testutils.CheckCommandFlag(cleanCmd, "keep", "", core.CfgKeyXCleanKeep, "3", false)
		testutils.CheckCommandFlag(cleanCmd, "name", "n", core.CfgKeyXCleanName, "", false)
	})

	Describe("helpers", func() {
		It("should return linux trash dir", func() {
			dir, err := defaultCleanTrashDir()
			Expect(err).NotTo(HaveOccurred())
			Expect(dir).To(Equal(filepath.Join(string(os.PathSeparator), "tmp", "cmdr-cleaned")))
		})

		It("should ensure dir and choose a unique path", func() {
			dir, err := os.MkdirTemp("", "cmdr-clean")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(dir)

			Expect(ensureDir(filepath.Join(dir, "nested"))).To(Succeed())
			path, err := uniquePath(dir, "new")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(dir, "new")))

			Expect(os.WriteFile(filepath.Join(dir, "cmd"), []byte("old"), 0644)).To(Succeed())

			path, err = uniquePath(dir, "cmd")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(dir, "cmd-1")))
		})

		It("should move files and report missing sources", func() {
			dir, err := os.MkdirTemp("", "cmdr-clean")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(dir)

			src := filepath.Join(dir, "src")
			dst := filepath.Join(dir, "dst")
			Expect(os.WriteFile(src, []byte("data"), 0600)).To(Succeed())
			Expect(moveFile(src, dst)).To(Succeed())
			Expect(dst).To(BeAnExistingFile())
			Expect(src).NotTo(BeAnExistingFile())
			Expect(moveFile(src, filepath.Join(dir, "missing"))).To(HaveOccurred())
		})

		It("should copy and remove files on cross-device rename", func() {
			dir, err := os.MkdirTemp("", "cmdr-clean")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(dir)

			src := filepath.Join(dir, "src")
			dst := filepath.Join(dir, "dst")
			Expect(os.WriteFile(src, []byte("data"), 0600)).To(Succeed())
			err = moveFileWithRename(src, dst, func(oldname, newname string) error {
				return &os.LinkError{Op: "rename", Old: oldname, New: newname, Err: syscall.EXDEV}
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(dst).To(BeAnExistingFile())
			Expect(src).NotTo(BeAnExistingFile())
			content, err := os.ReadFile(dst)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("data"))
		})

		It("should return cross-device fallback errors", func() {
			dir, err := os.MkdirTemp("", "cmdr-clean")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(dir)

			linkErr := func(oldname, newname string) error {
				return &os.LinkError{Op: "rename", Old: oldname, New: newname, Err: syscall.EXDEV}
			}
			Expect(moveFileWithRename("missing", filepath.Join(dir, "dst"), linkErr)).To(HaveOccurred())

			src := filepath.Join(dir, "src")
			Expect(os.WriteFile(src, []byte("data"), 0600)).To(Succeed())
			Expect(moveFileWithRename(src, filepath.Join(dir, "missing", "dst"), linkErr)).To(HaveOccurred())

			srcDir := filepath.Join(dir, "src-dir")
			Expect(os.Mkdir(srcDir, 0755)).To(Succeed())
			Expect(moveFileWithRename(srcDir, filepath.Join(dir, "dst-dir"), linkErr)).To(HaveOccurred())

			Expect(moveFileWithRename(src, filepath.Join(dir, "dst"), func(oldname, newname string) error {
				return errors.New("rename failed")
			})).To(HaveOccurred())
		})
	})

	Describe("runClean", func() {
		var (
			cfg      *viper.Viper
			root     string
			binDir   string
			trashDir string
			now      time.Time
			deps     cleanDeps
		)

		newShim := func(name, version string, age time.Duration) string {
			path := filepath.Join(root, "shims", name+"-"+version)
			Expect(os.WriteFile(path, []byte(version), 0755)).To(Succeed())
			ts := now.Add(-age)
			Expect(os.Chtimes(path, ts, ts)).To(Succeed())
			return path
		}

		BeforeEach(func() {
			var err error
			root, err = os.MkdirTemp("", "cmdr-clean")
			Expect(err).NotTo(HaveOccurred())
			binDir = filepath.Join(root, "bin")
			trashDir = filepath.Join(root, "trash")
			Expect(os.MkdirAll(binDir, 0755)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(root, "shims"), 0755)).To(Succeed())

			now = time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
			cfg = viper.New()
			cfg.Set(core.CfgKeyCmdrBinDir, binDir)
			cfg.Set(core.CfgKeyXCleanAgeDays, 30)
			cfg.Set(core.CfgKeyXCleanKeep, 1)
			cfg.Set(core.CfgKeyXCleanName, []string{})
			deps = defaultCleanDeps()
			deps.now = func() time.Time { return now }
			deps.trashDir = func() (string, error) { return trashDir, nil }
		})

		AfterEach(func() {
			os.RemoveAll(root)
		})

		It("should reject invalid age and keep", func() {
			manager := &cleanTestManager{}
			cfg.Set(core.CfgKeyXCleanAgeDays, -1)
			Expect(runClean(cfg, manager, deps)).To(MatchError("age must be >= 0"))

			cfg.Set(core.CfgKeyXCleanAgeDays, 1)
			cfg.Set(core.CfgKeyXCleanKeep, -1)
			Expect(runClean(cfg, manager, deps)).To(MatchError("keep must be >= 0"))
		})

		It("should validate wanted names", func() {
			manager := &cleanTestManager{commands: []core.Command{
				cleanTestCommand{name: "cmdr", version: "1.0.0", location: newShim("cmdr", "1.0.0", 60*24*time.Hour)},
			}}

			cfg.Set(core.CfgKeyXCleanName, []string{"  "})
			Expect(runClean(cfg, manager, deps)).To(MatchError("name must not be empty"))

			cfg.Set(core.CfgKeyXCleanName, []string{"missing"})
			Expect(runClean(cfg, manager, deps)).To(MatchError("command(s) not found: missing"))
		})

		It("should clean old inactive versions after keeping newest candidates", func() {
			newest := newShim("cmdr", "3.0.0", 90*24*time.Hour)
			oldest := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			recent := newShim("other", "1.0.0", 5*24*time.Hour)
			manager := &cleanTestManager{commands: []core.Command{
				cleanTestCommand{name: "cmdr", version: "3.0.0", location: newest},
				cleanTestCommand{name: "cmdr", version: "1.0.0", location: oldest},
				cleanTestCommand{name: "other", version: "1.0.0", location: recent},
			}}

			Expect(runClean(cfg, manager, deps)).To(Succeed())
			Expect(manager.undefined).To(Equal([]string{"cmdr:1.0.0"}))
			Expect(oldest).NotTo(BeAnExistingFile())
			Expect(filepath.Join(trashDir, "cmdr", filepath.Base(oldest))).To(BeAnExistingFile())
			Expect(newest).To(BeAnExistingFile())
			Expect(recent).To(BeAnExistingFile())
		})

		It("should filter wanted names and skip missing shim files", func() {
			target := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			other := newShim("other", "1.0.0", 120*24*time.Hour)
			missing := filepath.Join(root, "shims", "missing-1.0.0")
			manager := &cleanTestManager{commands: []core.Command{
				cleanTestCommand{name: "cmdr", version: "1.0.0", location: target},
				cleanTestCommand{name: "cmdr", version: "0.9.0", location: missing},
				cleanTestCommand{name: "other", version: "1.0.0", location: other},
			}}
			cfg.Set(core.CfgKeyXCleanKeep, 0)
			cfg.Set(core.CfgKeyXCleanName, []string{"cmdr"})

			Expect(runClean(cfg, manager, deps)).To(Succeed())
			Expect(manager.undefined).To(Equal([]string{"cmdr:1.0.0"}))
			Expect(target).NotTo(BeAnExistingFile())
			Expect(other).To(BeAnExistingFile())
		})

		It("should skip activated commands and shim-detected active locations", func() {
			active := newShim("cmdr", "2.0.0", 120*24*time.Hour)
			inactive := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			Expect(os.Symlink(inactive, filepath.Join(binDir, "cmdr"))).To(Succeed())
			manager := &cleanTestManager{commands: []core.Command{
				cleanTestCommand{name: "cmdr", version: "2.0.0", activated: true, location: active},
				cleanTestCommand{name: "cmdr", version: "1.0.0", location: inactive},
			}}
			cfg.Set(core.CfgKeyXCleanKeep, 0)

			Expect(runClean(cfg, manager, deps)).To(Succeed())
			Expect(manager.undefined).To(BeEmpty())
			Expect(active).To(BeAnExistingFile())
			Expect(inactive).To(BeAnExistingFile())
		})

		It("should return query, trash, move and undefine errors", func() {
			manager := &cleanTestManager{queryErr: errors.New("query failed")}
			Expect(runClean(cfg, manager, deps)).To(MatchError("query failed"))

			manager = &cleanTestManager{allErr: errors.New("all failed")}
			Expect(runClean(cfg, manager, deps)).To(MatchError("all failed"))

			deps.trashDir = func() (string, error) { return "", errors.New("trash failed") }
			Expect(runClean(cfg, manager, deps)).To(MatchError("trash failed"))

			deps = defaultCleanDeps()
			deps.now = func() time.Time { return now }
			deps.trashDir = func() (string, error) { return trashDir, nil }
			old := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			manager = &cleanTestManager{
				commands: []core.Command{cleanTestCommand{name: "cmdr", version: "1.0.0", location: old}},
				undefErr: errors.New("undefine failed"),
			}
			cfg.Set(core.CfgKeyXCleanKeep, 0)
			Expect(runClean(cfg, manager, deps)).To(HaveOccurred())
			Expect(old).To(BeAnExistingFile())

			deps.move = func(string, string) error { return errors.New("move failed") }
			Expect(runClean(cfg, manager, deps)).To(HaveOccurred())
		})

		It("should report undefine errors even when rollback fails", func() {
			old := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			manager := &cleanTestManager{
				commands: []core.Command{cleanTestCommand{name: "cmdr", version: "1.0.0", location: old}},
				undefErr: errors.New("undefine failed"),
			}
			cfg.Set(core.CfgKeyXCleanKeep, 0)

			moved := false
			deps.move = func(src, dst string) error {
				if !moved {
					moved = true
					return os.Rename(src, dst)
				}
				return errors.New("rollback failed")
			}

			err := runClean(cfg, manager, deps)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("undefine cmdr:1.0.0 failed"))
		})

		It("should collect destination preparation errors", func() {
			old := newShim("cmdr", "1.0.0", 120*24*time.Hour)
			manager := &cleanTestManager{commands: []core.Command{
				cleanTestCommand{name: "cmdr", version: "1.0.0", location: old},
			}}
			cfg.Set(core.CfgKeyXCleanKeep, 0)

			deps.ensure = func(path string) error {
				if path == trashDir {
					return nil
				}
				return errors.New("ensure failed")
			}
			Expect(runClean(cfg, manager, deps)).To(HaveOccurred())

			deps = defaultCleanDeps()
			deps.now = func() time.Time { return now }
			deps.trashDir = func() (string, error) { return trashDir, nil }
			deps.unique = func(string, string) (string, error) { return "", errors.New("unique failed") }
			Expect(runClean(cfg, manager, deps)).To(HaveOccurred())
		})
	})
})
