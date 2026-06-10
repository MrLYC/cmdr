package strategy

import (
	"errors"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

type temporaryNetError struct{}

func (temporaryNetError) Error() string   { return "temporary" }
func (temporaryNetError) Timeout() bool   { return false }
func (temporaryNetError) Temporary() bool { return true }

var _ = Describe("Strategy misc", func() {
	It("should validate strategy config", func() {
		Expect((&StrategyConfig{Timeout: -1}).Validate()).To(MatchError("invalid timeout: -1"))
		Expect((&StrategyConfig{MaxRetries: -1}).Validate()).To(MatchError("invalid max retries: -1"))
		Expect((&StrategyConfig{EnableProxy: true}).Validate()).To(MatchError("proxy type is required when proxy is enabled"))
		Expect((&StrategyConfig{EnableProxy: true, ProxyType: "ftp"}).Validate()).To(MatchError("unsupported proxy type: ftp (supported: http, socks5)"))
		Expect((&StrategyConfig{EnableProxy: true, ProxyType: "http"}).Validate()).To(MatchError("proxy address is required when proxy is enabled"))
		Expect((&StrategyConfig{EnableProxy: true, ProxyType: "http", ProxyAddr: "://bad"}).Validate()).To(HaveOccurred())
		Expect((&StrategyConfig{EnableProxy: true, ProxyType: "http", ProxyAddr: "http://127.0.0.1:8080"}).Validate()).To(Succeed())
	})

	It("should expose strategy chain helpers", func() {
		cfg := viper.New()
		cfg.Set("download.direct.max_retries", 1)
		cfg.Set("download.proxy.enabled", false)

		chain := NewStrategyChain()
		Expect(chain.Strategies()).To(BeEmpty())
		direct := NewDirectStrategy()
		chain.AddStrategy(direct)
		chain.AddStrategy(NewProxyStrategy())
		Expect(chain.Strategies()).To(HaveLen(2))
		Expect(chain.Configure(cfg)).To(Succeed())
		Expect(chain.Execute("https://example.com/file", func(string) error { return nil })).To(Succeed())
	})

	It("should cover direct and proxy retry decisions", func() {
		direct := NewDirectStrategy()
		direct.SetEnabled(true)
		Expect(direct.IsEnabled("https://example.com")).To(BeTrue())
		Expect(direct.ShouldRetry(nil)).To(BeFalse())
		Expect(direct.ShouldRetry(errors.New("deadline exceeded"))).To(BeTrue())
		Expect(direct.ShouldRetry(errors.New("connection reset"))).To(BeTrue())
		var netErr net.Error = temporaryNetError{}
		Expect(direct.ShouldRetry(netErr)).To(BeTrue())
		Expect(direct.ShouldFallback(nil)).To(BeFalse())
		Expect(direct.ShouldFallback(ErrNetworkError)).To(BeTrue())

		proxy := NewProxyStrategy()
		proxy.SetEnabled(true)
		Expect(proxy.IsEnabled("https://example.com")).To(BeTrue())
		Expect(proxy.ShouldRetry(errors.New("dial tcp failed"))).To(BeTrue())
		Expect(proxy.ShouldFallback(ErrConnectionError)).To(BeTrue())
		Expect(proxy.GetOptions()).To(BeNil())
	})

	It("should rewrite URIs and expose rewrite state", func() {
		cfg := viper.New()
		cfg.Set("download.rewrite.rule", "{{ .Scheme }}://mirror.local{{ .Path }}?{{ .Query }}#{{ .Fragment }}")
		cfg.Set("download.rewrite.condition.schemes", []string{"https"})
		cfg.Set("download.rewrite.condition.hosts", []string{"github.com"})
		cfg.Set("download.rewrite.condition.patterns", []string{"github.*"})

		rewrite := NewRewriteStrategy()
		Expect(rewrite.Configure(cfg)).To(Succeed())
		Expect(rewrite.IsEnabledConfigured()).To(BeTrue())
		Expect(rewrite.IsEnabled("https://github.com/mrlyc/cmdr?download=1#asset")).To(BeTrue())
		Expect(rewrite.ShouldRetry(errors.New("anything"))).To(BeFalse())
		Expect(rewrite.ShouldFallback(errors.New("anything"))).To(BeFalse())

		uri, err := rewrite.GetRewrittenURI("https://github.com/mrlyc/cmdr?download=1#asset")
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(Equal("https://mirror.local/mrlyc/cmdr?download=1#asset"))

		rewrite.SetEnabled(false)
		Expect(rewrite.IsEnabled("https://github.com/mrlyc/cmdr")).To(BeTrue())

		cfg = viper.New()
		cfg.Set("download.rewrite.rule", "{{")
		Expect(NewRewriteStrategy().Configure(cfg)).To(HaveOccurred())
	})
})
