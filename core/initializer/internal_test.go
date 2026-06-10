package initializer

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
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
	queryErr  error
	query     core.CommandQuery
	defined   bool
	activated bool
	undefined []string
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

func TestProfilePathByShell(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"bash": filepath.Join(home, ".bashrc"),
		"zsh":  filepath.Join(home, ".zshrc"),
		"fish": filepath.Join(home, ".config", "fish", "config.fish"),
		"sh":   filepath.Join(home, ".profile"),
		"ash":  filepath.Join(home, ".profile"),
	}

	for shell, expected := range tests {
		got, err := getProfilePathByShell(shell)
		if err != nil {
			t.Fatalf("getProfilePathByShell(%s): %v", shell, err)
		}
		if got != expected {
			t.Fatalf("getProfilePathByShell(%s) = %s, want %s", shell, got, expected)
		}
	}
}

func TestInitializerFactories(t *testing.T) {
	dir := t.TempDir()
	cfg := viper.New()
	cfg.Set(core.CfgKeyCmdrProfileDir, dir)
	cfg.Set(core.CfgKeyCmdrProfilePath, filepath.Join(dir, "profile"))
	cfg.Set(core.CfgKeyCmdrShell, "sh")

	for _, key := range []string{"profile-dir-backup", "profile-dir-export", "profile-dir-render", "profile-injector", "database-migrator"} {
		initializer, err := core.NewInitializer(key, cfg)
		if err != nil {
			t.Fatalf("NewInitializer(%s): %v", key, err)
		}
		if initializer == nil {
			t.Fatalf("NewInitializer(%s) returned nil", key)
		}
	}
}

func TestCmdrUpdaterFactory(t *testing.T) {
	previous := core.GetCommandManagerFactory(core.CommandProviderDatabase)
	defer core.RegisterCommandManagerFactory(core.CommandProviderDatabase, previous)

	core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(cfg core.Configuration) (core.CommandManager, error) {
		return initializerCommandManager{query: initializerCommandQuery{count: 0}}, nil
	})

	initializer, err := core.NewInitializer("cmdr-updater", viper.New())
	if err != nil {
		t.Fatal(err)
	}
	if initializer == nil {
		t.Fatal("expected initializer")
	}
}

func TestEmbedFSExporterErrorBranches(t *testing.T) {
	root := t.TempDir()
	dst := t.TempDir()
	exporter := NewEmbedFSExporter(os.DirFS(root), "root", dst, 0644)

	blocker := filepath.Join(dst, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := exporter.copyDir(blocker, 0755); err == nil {
		t.Fatal("expected copyDir error")
	}
	if err := exporter.copyFile("missing", filepath.Join(dst, "out"), 0644); err == nil {
		t.Fatal("expected source open error")
	}

	src := filepath.Join(root, "source")
	if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := exporter.copyFile("source", filepath.Join(dst, "missing", "out"), 0644); err == nil {
		t.Fatal("expected destination open error")
	}
	if err := exporter.exportDir("root/file", nil, errors.New("walk failed")); err == nil {
		t.Fatal("expected exportDir walk error")
	}
	if err := exporter.Init(false); err == nil {
		t.Fatal("expected missing source error")
	}
}

func TestDirRenderErrorBranches(t *testing.T) {
	root := t.TempDir()
	renderer := NewDirRender(root, ".gotmpl", map[string]string{"key": "value"})

	if err := renderer.renderTemplate("bad", "{{", bytes.NewBuffer(nil)); err == nil {
		t.Fatal("expected template parse error")
	}
	if err := renderer.renderTemplate("write", "content", failingWriter{}); err == nil {
		t.Fatal("expected template execute/write error")
	}

	src := filepath.Join(root, "source.txt.gotmpl")
	if err := os.WriteFile(src, []byte("{{ .key }}"), 0644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "source.txt")
	if err := os.Mkdir(target, 0755); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := renderer.renderFile(src, info); err == nil {
		t.Fatal("expected renderFile create target error")
	}

	dirTemplate := filepath.Join(root, "blocker", "child.gotmpl")
	if err := os.WriteFile(filepath.Join(root, "blocker"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := renderer.renderDir(dirTemplate, info); err == nil {
		t.Fatal("expected renderDir create target error")
	}
}

func TestProfileInjectorErrorBranches(t *testing.T) {
	root := t.TempDir()
	missingProfile := filepath.Join(root, "missing")
	injector := NewProfileInjector("/tmp/cmdr_initializer.sh", missingProfile)
	if _, err := injector.makeProfileScript(); err == nil {
		t.Fatal("expected makeProfileScript open error")
	}

	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	injector = NewProfileInjector("/tmp/cmdr_initializer.sh", filepath.Join(blocker, "profile"))
	if err := injector.Init(false); err == nil {
		t.Fatal("expected profile open error")
	}
}

func TestCmdrUpdaterCollectLegacyVersionErrors(t *testing.T) {
	tests := []core.CommandManager{
		initializerCommandManager{queryErr: errors.New("query failed")},
		initializerCommandManager{query: initializerCommandQuery{countErr: errors.New("count failed")}},
		initializerCommandManager{query: initializerCommandQuery{allErr: errors.New("all failed")}},
	}

	for _, manager := range tests {
		updater := NewCmdrUpdater(manager, "cmdr", "1.0.0", "/tmp/cmdr")
		if err := updater.Init(true); err == nil {
			t.Fatal("expected error")
		}
	}
}

func TestCmdrUpdaterCollectLegacyVersionBranches(t *testing.T) {
	updater := NewCmdrUpdater(
		initializerCommandManager{query: initializerCommandQuery{count: 0}},
		core.Name,
		"2.0.0",
		"/tmp/cmdr",
	)
	if err := updater.Init(true); err != nil {
		t.Fatal(err)
	}

	updater = NewCmdrUpdater(
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
	if err := updater.Init(true); err != nil {
		t.Fatal(err)
	}
}
