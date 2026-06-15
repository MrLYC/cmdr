package manager

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("Download manager internals", func() {
	It("keeps prepared URI when direct strategy is selected", func() {
		manager := NewDownloadManager(nil, nil, 0, nil)
		manager.SetStrategyChain(strategy.NewStrategyChain(strategy.NewDirectStrategy()))

		Expect(manager.maybeRewriteURI("prepared", "original")).To(Equal("prepared"))
	})

	It("returns original URI when rewrite keeps it unchanged", func() {
		cfg := viper.New()
		cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .URI }}")
		rewrite := strategy.NewRewriteStrategy()
		Expect(rewrite.Configure(cfg)).To(Succeed())

		manager := NewDownloadManager(nil, nil, 0, nil)
		manager.SetStrategyChain(strategy.NewStrategyChain(rewrite))

		Expect(manager.maybeRewriteURI("https://example.com/cmd", "https://example.com/cmd")).To(Equal("https://example.com/cmd"))
	})

	It("resets supported fetcher option types", func() {
		resetFetcherOptions(noopFetcher{})
		resetFetcherOptions(fetcher.NewDefaultGoGetter(io.Discard))
	})
})
