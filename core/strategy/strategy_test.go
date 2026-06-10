package strategy

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
)

type testStrategy struct {
	name           string
	enabled        bool
	preparedURI    string
	prepareErr     error
	shouldRetry    bool
	shouldFallback bool
	configureErr   error
}

func (s *testStrategy) Name() string { return s.name }
func (s *testStrategy) Prepare(uri string) (string, error) {
	if s.prepareErr != nil {
		return "", s.prepareErr
	}
	if s.preparedURI != "" {
		return s.preparedURI, nil
	}
	return uri, nil
}
func (s *testStrategy) ShouldRetry(error) bool    { return s.shouldRetry }
func (s *testStrategy) ShouldFallback(error) bool { return s.shouldFallback }
func (s *testStrategy) Configure(core.Configuration) error {
	return s.configureErr
}
func (s *testStrategy) IsEnabled(string) bool   { return s.enabled }
func (s *testStrategy) SetEnabled(enabled bool) { s.enabled = enabled }

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

		// GetOptions() returns nil because go-getter v1.8.4 doesn't support
		// WithTimeout/WithProxy as ClientOption. The strategy uses other
		// mechanisms for configuration.
		options := strategy.GetOptions()
		Expect(options).To(BeNil())
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

		// GetOptions() returns nil - proxy settings are applied via other mechanisms
		options := strategy.GetOptions()
		Expect(options).To(BeNil())
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
		proxy := NewProxyStrategy()
		cfg := viper.New()
		direct.Configure(cfg)
		cfg.Set("download.proxy.enabled", true)
		cfg.Set("download.proxy.address", "http://proxy:8080")
		proxy.Configure(cfg)

		chain := NewStrategyChain(direct, proxy)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			attemptCount++
			return ErrNetworkError
		})

		Expect(err).NotTo(BeNil())
		Expect(attemptCount).To(BeNumerically(">", 1))
	})

	It("should use only enabled strategies", func() {
		enabledCount := 0

		direct := NewDirectStrategy()
		proxy := NewProxyStrategy()
		cfg := viper.New()

		direct.Configure(cfg)
		cfg.Set("download.proxy.enabled", false)
		proxy.Configure(cfg)

		chain := NewStrategyChain(direct, proxy)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			enabledCount++
			return nil
		})

		Expect(err).To(BeNil())
		Expect(enabledCount).To(Equal(1))
	})

	It("should skip strategies that fail prepare", func() {
		calls := []string{}
		first := &testStrategy{name: "first", enabled: true, prepareErr: errors.New("bad uri")}
		second := &testStrategy{name: "second", enabled: true, preparedURI: "prepared"}
		chain := NewStrategyChain(first, second)

		err := chain.Execute("https://example.com/file", func(uri string) error {
			calls = append(calls, uri)
			return nil
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(calls).To(Equal([]string{"prepared"}))
	})

	It("should retry the same strategy when requested", func() {
		attempts := 0
		chain := NewStrategyChain(&testStrategy{name: "retry", enabled: true, shouldRetry: true})

		err := chain.Execute("https://example.com/file", func(uri string) error {
			attempts++
			if attempts == 1 {
				return errors.New("temporary")
			}
			return nil
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(attempts).To(Equal(2))
	})

	It("should return non-retriable strategy errors directly", func() {
		expected := errors.New("permanent")
		chain := NewStrategyChain(&testStrategy{name: "direct", enabled: true})

		err := chain.Execute("https://example.com/file", func(uri string) error {
			return expected
		})

		Expect(err).To(MatchError(expected))
	})

	It("should report unexpected state when no strategy runs", func() {
		chain := NewStrategyChain(&testStrategy{name: "bad", enabled: true, prepareErr: errors.New("prepare failed")})
		Expect(chain.Execute("https://example.com/file", func(uri string) error {
			return nil
		})).To(MatchError("unexpected state: no error but download failed"))
	})

	It("should configure strategy errors with strategy names", func() {
		chain := NewStrategyChain(&testStrategy{name: "bad", configureErr: errors.New("invalid")})
		Expect(chain.Configure(viper.New())).To(MatchError(ContainSubstring("failed to configure strategy bad")))
	})
})
