package manager

import (
	"io"
	"testing"

	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/fetcher"
	"github.com/mrlyc/cmdr/core/strategy"
)

type noopFetcher struct{}

func (noopFetcher) IsSupport(string) bool {
	return true
}

func (noopFetcher) Fetch(string, string, string, string) error {
	return nil
}

func TestDownloadManagerMaybeRewriteURI(t *testing.T) {
	manager := NewDownloadManager(nil, nil, 0, nil)
	manager.SetStrategyChain(strategy.NewStrategyChain(strategy.NewDirectStrategy()))
	if got := manager.maybeRewriteURI("prepared", "original"); got != "prepared" {
		t.Fatalf("expected prepared URI, got %s", got)
	}

	cfg := viper.New()
	cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .URI }}")
	rewrite := strategy.NewRewriteStrategy()
	if err := rewrite.Configure(cfg); err != nil {
		t.Fatal(err)
	}
	manager.SetStrategyChain(strategy.NewStrategyChain(rewrite))
	if got := manager.maybeRewriteURI("https://example.com/cmd", "https://example.com/cmd"); got != "https://example.com/cmd" {
		t.Fatalf("expected original URI, got %s", got)
	}
}

func TestResetFetcherOptions(t *testing.T) {
	resetFetcherOptions(noopFetcher{})
	resetFetcherOptions(fetcher.NewDefaultGoGetter(io.Discard))
}
