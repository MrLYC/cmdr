package initializer

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/internal/testutils"
)

type failingWriter struct{}

func (f failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

type initializerCommandQuery struct {
	count    int
	countErr error
	allErr   error
	commands []core.Command
}

func (q initializerCommandQuery) WithName(string) core.CommandQuery     { return q }
func (q initializerCommandQuery) WithVersion(string) core.CommandQuery  { return q }
func (q initializerCommandQuery) WithActivated(bool) core.CommandQuery  { return q }
func (q initializerCommandQuery) WithLocation(string) core.CommandQuery { return q }
func (q initializerCommandQuery) All() ([]core.Command, error) {
	if q.allErr != nil {
		return nil, q.allErr
	}
	return q.commands, nil
}
func (q initializerCommandQuery) One() (core.Command, error) { return nil, nil }
func (q initializerCommandQuery) Count() (int, error) {
	if q.countErr != nil {
		return 0, q.countErr
	}
	if q.count != 0 {
		return q.count, nil
	}
	if q.commands != nil {
		return len(q.commands), nil
	}
	return 1, nil
}

type initializerCommand struct {
	version   string
	activated bool
}

func (c initializerCommand) GetName() string     { return core.Name }
func (c initializerCommand) GetVersion() string  { return c.version }
func (c initializerCommand) GetActivated() bool  { return c.activated }
func (c initializerCommand) GetLocation() string { return "/tmp/cmdr_" + c.version }

type initializerCommandManager struct {
	queryErr error
	query    core.CommandQuery
}

func (m initializerCommandManager) Close() error { return nil }
func (m initializerCommandManager) Provider() core.CommandProvider {
	return core.CommandProviderDatabase
}
func (m initializerCommandManager) Query() (core.CommandQuery, error) {
	return m.query, m.queryErr
}
func (m initializerCommandManager) Define(string, string, string) (core.Command, error) {
	return nil, nil
}
func (m initializerCommandManager) Undefine(string, string) error { return nil }
func (m initializerCommandManager) Activate(string, string) error { return nil }
func (m initializerCommandManager) Deactivate(string) error       { return nil }

func initializerTempDir() string {
	dir, err := os.MkdirTemp("", "cmdr-initializer-test-*")
	Expect(err).NotTo(HaveOccurred())
	return dir
}

var _ = Describe("Initializer internals", func() {
	DescribeTable("getProfilePathByShell",
		func(shell string, expected func(string) string) {
			home, err := os.UserHomeDir()
			Expect(err).NotTo(HaveOccurred())

			got, err := getProfilePathByShell(shell)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(expected(home)))
		},
		Entry("bash", "bash", func(home string) string { return filepath.Join(home, ".bashrc") }),
		Entry("zsh", "zsh", func(home string) string { return filepath.Join(home, ".zshrc") }),
		Entry("fish", "fish", func(home string) string { return filepath.Join(home, ".config", "fish", "config.fish") }),
		Entry("sh", "sh", func(home string) string { return filepath.Join(home, ".profile") }),
		Entry("ash", "ash", func(home string) string { return filepath.Join(home, ".profile") }),
	)

	It("creates configured initializer factories", func() {
		dir := initializerTempDir()
		defer os.RemoveAll(dir)

		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrProfileDir, dir)
		cfg.Set(core.CfgKeyCmdrProfilePath, filepath.Join(dir, "profile"))
		cfg.Set(core.CfgKeyCmdrShell, "sh")

		for _, key := range []string{"profile-dir-backup", "profile-dir-export", "profile-dir-render", "profile-injector", "database-migrator"} {
			initializer, err := core.NewInitializer(key, cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(initializer).NotTo(BeNil())
		}
	})

	It("creates the cmdr updater factory with a database manager", func() {
		restoreFactory := testutils.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(core.Configuration) (core.CommandManager, error) {
			return initializerCommandManager{query: initializerCommandQuery{count: 0}}, nil
		})
		defer restoreFactory()

		initializer, err := core.NewInitializer("cmdr-updater", viper.New())
		Expect(err).NotTo(HaveOccurred())
		Expect(initializer).NotTo(BeNil())
	})

	Context("EmbedFSExporter", func() {
		It("returns filesystem errors", func() {
			root := initializerTempDir()
			defer os.RemoveAll(root)
			dst := initializerTempDir()
			defer os.RemoveAll(dst)
			exporter := NewEmbedFSExporter(os.DirFS(root), "root", dst, 0644)

			blocker := filepath.Join(dst, "blocker")
			Expect(os.WriteFile(blocker, []byte("x"), 0644)).To(Succeed())
			Expect(exporter.copyDir(blocker, 0755)).To(HaveOccurred())
			Expect(exporter.copyFile("missing", filepath.Join(dst, "out"), 0644)).To(HaveOccurred())

			src := filepath.Join(root, "source")
			Expect(os.WriteFile(src, []byte("x"), 0644)).To(Succeed())
			Expect(exporter.copyFile("source", filepath.Join(dst, "missing", "out"), 0644)).To(HaveOccurred())
			Expect(exporter.exportDir("root/file", nil, errors.New("walk failed"))).To(HaveOccurred())
			Expect(exporter.Init(false)).To(HaveOccurred())
		})
	})

	Context("DirRender", func() {
		It("returns template and filesystem errors", func() {
			root := initializerTempDir()
			defer os.RemoveAll(root)
			renderer := NewDirRender(root, ".gotmpl", map[string]string{"key": "value"})

			Expect(renderer.renderTemplate("bad", "{{", bytes.NewBuffer(nil))).To(HaveOccurred())
			Expect(renderer.renderTemplate("write", "content", failingWriter{})).To(HaveOccurred())

			src := filepath.Join(root, "source.txt.gotmpl")
			Expect(os.WriteFile(src, []byte("{{ .key }}"), 0644)).To(Succeed())
			target := filepath.Join(root, "source.txt")
			Expect(os.Mkdir(target, 0755)).To(Succeed())
			info, err := os.Stat(src)
			Expect(err).NotTo(HaveOccurred())
			Expect(renderer.renderFile(src, info)).To(HaveOccurred())

			dirTemplate := filepath.Join(root, "blocker", "child.gotmpl")
			Expect(os.WriteFile(filepath.Join(root, "blocker"), []byte("x"), 0644)).To(Succeed())
			Expect(renderer.renderDir(dirTemplate, info)).To(HaveOccurred())
		})
	})

	Context("ProfileInjector", func() {
		It("returns profile read and write errors", func() {
			root := initializerTempDir()
			defer os.RemoveAll(root)

			missingProfile := filepath.Join(root, "missing")
			injector := NewProfileInjector("/tmp/cmdr_initializer.sh", missingProfile)
			_, err := injector.makeProfileScript()
			Expect(err).To(HaveOccurred())

			blocker := filepath.Join(root, "blocker")
			Expect(os.WriteFile(blocker, []byte("x"), 0644)).To(Succeed())
			injector = NewProfileInjector("/tmp/cmdr_initializer.sh", filepath.Join(blocker, "profile"))
			Expect(injector.Init(false)).To(HaveOccurred())
		})
	})

	Context("CmdrUpdater", func() {
		DescribeTable("returns legacy collection errors",
			func(manager core.CommandManager) {
				updater := NewCmdrUpdater(manager, "cmdr", "1.0.0", "/tmp/cmdr")
				Expect(updater.Init(true)).To(HaveOccurred())
			},
			Entry("query failure", initializerCommandManager{queryErr: errors.New("query failed")}),
			Entry("count failure", initializerCommandManager{query: initializerCommandQuery{countErr: errors.New("count failed")}}),
			Entry("all failure", initializerCommandManager{query: initializerCommandQuery{allErr: errors.New("all failed")}}),
		)

		It("handles empty legacy command sets", func() {
			updater := NewCmdrUpdater(
				initializerCommandManager{query: initializerCommandQuery{count: 0}},
				core.Name,
				"2.0.0",
				"/tmp/cmdr",
			)
			Expect(updater.Init(true)).To(Succeed())
		})

		It("handles multiple legacy command branches", func() {
			updater := NewCmdrUpdater(
				initializerCommandManager{query: initializerCommandQuery{commands: []core.Command{
					initializerCommand{version: "1.0.0"},
					initializerCommand{version: "2.0.0"},
					initializerCommand{version: "3.0.0"},
					initializerCommand{version: "0.9.0", activated: true},
				}}},
				core.Name,
				"2.0.0",
				"/tmp/cmdr",
			)
			Expect(updater.Init(true)).To(Succeed())
		})
	})
})
