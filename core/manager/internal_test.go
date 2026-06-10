package manager

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/strategy"
)

type internalFetcher struct {
	seenURI string
	err     error
}

func (f *internalFetcher) IsSupport(string) bool { return true }
func (f *internalFetcher) Fetch(name, version, uri, dst string) error {
	f.seenURI = uri
	if f.err != nil {
		return f.err
	}
	return os.WriteFile(filepath.Join(dst, name), []byte("binary"), 0755)
}

func TestResolveVersionFromLocation(t *testing.T) {
	tests := []struct {
		name     string
		fallback string
		location string
		expected string
	}{
		{"cmdr", "1", "/tmp/cmdr_1.2.3", "1.2.3"},
		{"cmdr", "1", "/tmp/other_1.2.3", "1"},
		{"cmdr", "1", "/tmp/cmdr_", "1"},
	}
	for _, tt := range tests {
		if got := resolveVersionFromLocation(tt.name, tt.fallback, tt.location); got != tt.expected {
			t.Fatalf("resolveVersionFromLocation() = %s, want %s", got, tt.expected)
		}
	}
}

func TestDownloadManagerInternals(t *testing.T) {
	manager := NewDownloadManager(nil, nil, 1, nil)
	if _, err := manager.search("cmdr", filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected search error")
	}

	cfg := viper.New()
	cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .URI }}?mirror=1")
	rewrite := strategy.NewRewriteStrategy()
	if err := rewrite.Configure(cfg); err != nil {
		t.Fatal(err)
	}
	fetcher := &internalFetcher{}
	manager.SetStrategyChain(strategy.NewStrategyChain(rewrite))
	dst := t.TempDir()
	result, err := manager.fetch(fetcher, "cmdr", "1.0.0", "https://example.com/cmdr", dst)
	if err != nil {
		t.Fatal(err)
	}
	if fetcher.seenURI != "https://example.com/cmdr?mirror=1" {
		t.Fatalf("seen uri = %s", fetcher.seenURI)
	}
	if result != filepath.Join(dst, "cmdr") {
		t.Fatalf("result = %s", result)
	}

	fetcher.err = errors.New("fetch failed")
	if _, err := manager.fetch(fetcher, "cmdr", "1.0.0", "https://example.com/cmdr", t.TempDir()); err == nil {
		t.Fatal("expected fetch error")
	}
}

func TestBinaryManagerInternals(t *testing.T) {
	cfg := viper.New()
	cfg.Set(core.CfgKeyCmdrBinDir, "bin")
	cfg.Set(core.CfgKeyCmdrShimsDir, "shims")
	cfg.Set(core.CfgKeyCmdrLinkMode, "link")
	if newBinaryManagerByConfiguration(cfg) == nil {
		t.Fatal("expected link manager")
	}
	cfg.Set(core.CfgKeyCmdrLinkMode, "copy")
	if NewBinaryManagerWithCopy("bin", "shims", 0755) == nil {
		t.Fatal("expected copy manager")
	}
	if newBinaryManagerByConfiguration(cfg) == nil {
		t.Fatal("expected copy manager from config")
	}

	mgr := NewBinaryManagerWithLink(filepath.Join(t.TempDir(), "bin"), filepath.Join(t.TempDir(), "missing"), 0755)
	if _, err := mgr.Query(); err == nil {
		t.Fatal("expected query error")
	}
}
