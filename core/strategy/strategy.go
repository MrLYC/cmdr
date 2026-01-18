package strategy

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mrlyc/cmdr/core"
)

var (
	ErrAllStrategiesFailed = errors.New("all download strategies failed")
)

type DownloadStrategy interface {
	// Name returns strategy name
	Name() string

	// Prepare prepares of download with a strategy
	// Returns to modified URI and error
	Prepare(uri string) (string, error)

	// ShouldRetry determines if error indicates we should retry with this strategy
	ShouldRetry(err error) bool

	// ShouldFallback determines if error indicates we should try to next strategy
	ShouldFallback(err error) bool

	// Configure configures strategy with a given configuration
	Configure(cfg core.Configuration) error

	// IsEnabled determines if strategy is enabled for a given URI
	IsEnabled(uri string) bool

	// SetEnabled enables or disables the strategy
	SetEnabled(enabled bool)
}

type StrategyCondition struct {
	Schemes  []string // http, https, git, etc.
	Hosts    []string // github.com, nodejs.org, etc.
	Patterns []string // glob patterns for matching hosts
}

type StrategyConfig struct {
	Enabled     bool
	Timeout     int
	MaxRetries  int
	EnableProxy bool
	ProxyType   string // "http" or "socks5"
	ProxyAddr   string
	RewriteRule string
	Condition   *StrategyCondition
}

func (c *StrategyConfig) Validate() error {
	if c.Timeout < 0 {
		return fmt.Errorf("invalid timeout: %d", c.Timeout)
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("invalid max retries: %d", c.MaxRetries)
	}
	if c.EnableProxy {
		if c.ProxyType == "" {
			return fmt.Errorf("proxy type is required when proxy is enabled")
		}
		if c.ProxyType != "http" && c.ProxyType != "socks5" {
			return fmt.Errorf("unsupported proxy type: %s (supported: http, socks5)", c.ProxyType)
		}
		if c.ProxyAddr == "" {
			return fmt.Errorf("proxy address is required when proxy is enabled")
		}
		if _, err := url.Parse(c.ProxyAddr); err != nil {
			return fmt.Errorf("invalid proxy address: %w", err)
		}
	}
	return nil
}

func (c *StrategyConfig) Matches(uri string) bool {
	if c.Condition == nil {
		return c.Enabled
	}

	// Parse URI to extract scheme and host
	parsed, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// Check scheme
	if len(c.Condition.Schemes) > 0 {
		schemeMatch := false
		for _, scheme := range c.Condition.Schemes {
			if parsed.Scheme == scheme {
				schemeMatch = true
				break
			}
		}
		if !schemeMatch {
			return false
		}
	}

	// Check host
	if len(c.Condition.Hosts) > 0 {
		hostMatch := false
		for _, host := range c.Condition.Hosts {
			if parsed.Host == host || strings.HasSuffix(parsed.Host, "."+host) {
				hostMatch = true
				break
			}
		}
		if !hostMatch {
			return false
		}
	}

	// Check patterns
	if len(c.Condition.Patterns) > 0 {
		patternMatch := false
		for _, pattern := range c.Condition.Patterns {
			if matched, _ := filepath.Match(pattern, parsed.Host); matched {
				patternMatch = true
				break
			}
		}
		if !patternMatch {
			return false
		}
	}

	return true
}
