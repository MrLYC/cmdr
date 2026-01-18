package strategy

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/go-getter"
	"github.com/mrlyc/cmdr/core"
)

type ProxyStrategy struct {
	config   *StrategyConfig
	proxyURL *url.URL
	enabled  bool
}

func (s *ProxyStrategy) Name() string {
	return "proxy"
}

func (s *ProxyStrategy) Prepare(uri string) (string, error) {
	return uri, nil
}

func (s *ProxyStrategy) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	// Retry on timeout and connection errors
	return isTimeoutError(err) || isConnectionError(err)
}

func (s *ProxyStrategy) ShouldFallback(err error) bool {
	if err == nil {
		return false
	}

	// Fallback on network errors
	return isNetworkError(err)
}

func (s *ProxyStrategy) Configure(cfg core.Configuration) error {
	enableProxy := cfg.GetBool("download.proxy.enabled")
	if !enableProxy {
		s.config = &StrategyConfig{Enabled: false}
		s.enabled = false
		return nil
	}

	proxyType := cfg.GetString("download.proxy.type")
	proxyAddr := cfg.GetString("download.proxy.address")

	s.config = &StrategyConfig{
		Enabled:     true,
		Timeout:     cfg.GetInt("download.proxy.timeout"),
		MaxRetries:  cfg.GetInt("download.proxy.max_retries"),
		EnableProxy: true,
		ProxyType:   proxyType,
		ProxyAddr:   proxyAddr,
	}

	// Parse condition
	if cfg.IsSet("download.proxy.condition") {
		s.config.Condition = &StrategyCondition{}

		// Parse schemes
		schemes := cfg.GetStringSlice("download.proxy.condition.schemes")
		if len(schemes) > 0 {
			s.config.Condition.Schemes = schemes
		}

		// Parse hosts
		hosts := cfg.GetStringSlice("download.proxy.condition.hosts")
		if len(hosts) > 0 {
			s.config.Condition.Hosts = hosts
		}

		// Parse patterns
		patterns := cfg.GetStringSlice("download.proxy.condition.patterns")
		if len(patterns) > 0 {
			s.config.Condition.Patterns = patterns
		}

		// If condition is set, strategy is only enabled when matching
		if s.config.Condition.Schemes != nil || s.config.Condition.Hosts != nil || s.config.Condition.Patterns != nil {
			s.config.Enabled = false
		}
	}

	if s.config.Timeout == 0 {
		s.config.Timeout = 30
	}
	if s.config.MaxRetries == 0 {
		s.config.MaxRetries = 3
	}

	// Parse proxy URL
	parsedURL, err := url.Parse(proxyAddr)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	s.proxyURL = parsedURL
	s.enabled = s.config.Enabled

	return nil
}

func (s *ProxyStrategy) GetOptions() []getter.ClientOption {
	return nil
}

func (s *ProxyStrategy) IsEnabled(uri string) bool {
	if s.config == nil {
		return s.enabled
	}
	if s.config.Condition == nil {
		return s.enabled
	}
	return s.config.Matches(uri)
}

func (s *ProxyStrategy) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func NewProxyStrategy() *ProxyStrategy {
	return &ProxyStrategy{}
}
