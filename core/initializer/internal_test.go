package initializer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
)

type initializerCommandQuery struct {
	countErr error
	allErr   error
}

func (q initializerCommandQuery) WithName(string) core.CommandQuery     { return q }
func (q initializerCommandQuery) WithVersion(string) core.CommandQuery  { return q }
func (q initializerCommandQuery) WithActivated(bool) core.CommandQuery  { return q }
func (q initializerCommandQuery) WithLocation(string) core.CommandQuery { return q }
func (q initializerCommandQuery) All() ([]core.Command, error)          { return nil, q.allErr }
func (q initializerCommandQuery) One() (core.Command, error)            { return nil, nil }
func (q initializerCommandQuery) Count() (int, error) {
	if q.countErr != nil {
		return 0, q.countErr
	}
	return 1, nil
}

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
