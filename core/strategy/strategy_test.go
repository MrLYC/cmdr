package strategy

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
)

func TestStrategy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Strategy Suite")
}

var _ = Describe("DirectStrategy", func() {
	var (
		cfg      core.Configuration
		strategy *DirectStrategy
	)

	BeforeEach(func() {
		cfg = viper.New()
		strategy = NewDirectStrategy()
	})

	It("should have correct name", func() {
		Expect(strategy.Name()).To(Equal("direct"))
	})

	It("should not modify URI", func() {
		uri := "https://example.com/file.tar.gz"
		result, err := strategy.Prepare(uri)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(uri))
	})

	It("should configure with default values", func() {
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		options := strategy.GetOptions()
		Expect(options).NotTo(BeNil())
		Expect(len(options)).To(BeNumerically(">", 0))
	})

	It("should be enabled by default", func() {
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())
		Expect(strategy.IsEnabled("https://example.com/file")).To(BeTrue())
	})

	It("should be enabled when condition matches scheme", func() {
		cfg.Set("download.direct.condition.schemes", []string{"https"})
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://example.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("http://example.com/file")).To(BeFalse())
	})

	It("should be enabled when condition matches host", func() {
		cfg.Set("download.direct.condition.hosts", []string{"github.com"})
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("https://gitlab.com/file")).To(BeFalse())
	})

	It("should be enabled when condition matches pattern", func() {
		cfg.Set("download.direct.condition.patterns", []string{"*.github.com"})
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://api.github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("https://example.com/file")).To(BeFalse())
	})
})

var _ = Describe("ProxyStrategy", func() {
	var (
		cfg      core.Configuration
		strategy *ProxyStrategy
	)

	BeforeEach(func() {
		cfg = viper.New()
		strategy = NewProxyStrategy()
	})

	It("should have correct name", func() {
		Expect(strategy.Name()).To(Equal("proxy"))
	})

	It("should not modify URI", func() {
		uri := "https://example.com/file.tar.gz"
		result, err := strategy.Prepare(uri)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(uri))
	})

	It("should be disabled by default", func() {
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())
		Expect(strategy.IsEnabled("https://example.com/file")).To(BeFalse())
	})

	It("should configure with proxy settings", func() {
		cfg.Set("download.proxy.enabled", true)
		cfg.Set("download.proxy.type", "http")
		cfg.Set("download.proxy.address", "http://proxy.example.com:8080")
		cfg.Set("download.proxy.timeout", 60)
		cfg.Set("download.proxy.max_retries", 5)

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())
		Expect(strategy.IsEnabled("https://example.com/file")).To(BeTrue())

		options := strategy.GetOptions()
		Expect(options).NotTo(BeNil())
	})

	It("should be enabled only for configured schemes", func() {
		cfg.Set("download.proxy.enabled", true)
		cfg.Set("download.proxy.type", "http")
		cfg.Set("download.proxy.address", "http://proxy:8080")
		cfg.Set("download.proxy.condition.schemes", []string{"https"})

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("http://github.com/file")).To(BeFalse())
	})

	It("should be enabled only for configured hosts", func() {
		cfg.Set("download.proxy.enabled", true)
		cfg.Set("download.proxy.type", "http")
		cfg.Set("download.proxy.address", "http://proxy:8080")
		cfg.Set("download.proxy.condition.hosts", []string{"github.com"})

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("https://gitlab.com/file")).To(BeFalse())
	})
})

var _ = Describe("RewriteStrategy", func() {
	var (
		cfg      core.Configuration
		strategy *RewriteStrategy
	)

	BeforeEach(func() {
		cfg = viper.New()
		strategy = NewRewriteStrategy()
	})

	It("should have correct name", func() {
		Expect(strategy.Name()).To(Equal("rewrite"))
	})

	It("should not modify URI when not configured", func() {
		uri := "https://github.com/user/repo/archive/v1.0.tar.gz"
		result, err := strategy.Prepare(uri)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(uri))
	})

	It("should be disabled by default", func() {
		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())
		Expect(strategy.IsEnabled("https://example.com/file")).To(BeFalse())
	})

	It("should rewrite URI with template", func() {
		cfg.Set("download.rewrite.rule", "https://mirror.example.com/{{.Path}}")

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		uri := "https://github.com/user/repo/archive/v1.0.tar.gz"
		result, err := strategy.Prepare(uri)
		Expect(err).To(BeNil())
		Expect(result).To(ContainSubstring("mirror.example.com"))
		Expect(strategy.IsEnabled(uri)).To(BeTrue())
	})

	It("should be enabled only for configured schemes", func() {
		cfg.Set("download.rewrite.rule", "https://mirror.com{{.Path}}")
		cfg.Set("download.rewrite.condition.schemes", []string{"https"})

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("http://github.com/file")).To(BeFalse())
	})

	It("should be enabled only for configured hosts", func() {
		cfg.Set("download.rewrite.rule", "https://mirror.com{{.Path}}")
		cfg.Set("download.rewrite.condition.hosts", []string{"github.com"})

		err := strategy.Configure(cfg)
		Expect(err).To(BeNil())

		Expect(strategy.IsEnabled("https://github.com/file")).To(BeTrue())
		Expect(strategy.IsEnabled("https://gitlab.com/file")).To(BeFalse())
	})
})

var _ = Describe("StrategyChain", func() {
	It("should execute strategies in order", func() {
		executedOrder := []string{}

		direct := NewDirectStrategy()
		proxy := NewProxyStrategy()

		chain := NewStrategyChain(direct, proxy)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			executedOrder = append(executedOrder, uri)
			return nil
		})

		Expect(err).To(BeNil())
		Expect(len(executedOrder)).To(Equal(1))
	})

	It("should fallback to next strategy on error", func() {
		attemptCount := 0

		direct := NewDirectStrategy()

		chain := NewStrategyChain(direct)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			attemptCount++
			if attemptCount == 1 {
				return ErrNetworkError
			}
			return nil
		})

		Expect(err).NotTo(BeNil())
		Expect(attemptCount).To(BeNumerically(">", 1))
	})

	It("should use only enabled strategies", func() {
		enabledCount := 0

		proxy := NewProxyStrategy()
		cfg := viper.New()
		cfg.Set("download.proxy.enabled", false)

		proxy.Configure(cfg)

		chain := NewStrategyChain(proxy)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			enabledCount++
			return ErrNetworkError
		})

		Expect(enabledCount).To(Equal(0))
		Expect(err).NotTo(BeNil())
	})
})
