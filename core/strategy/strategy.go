package strategy

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/mrlyc/cmdr/core"
)

var (
	ErrAllStrategiesFailed = errors.New("all download strategies failed")
)

type DownloadStrategy interface {
	// Name returns the strategy name
	Name() string

	// Prepare prepares the download with the strategy
	// Returns the modified URI and error
	Prepare(uri string) (string, error)

	// ShouldRetry determines if the error indicates we should retry with this strategy
	ShouldRetry(err error) bool

	// ShouldFallback determines if the error indicates we should try the next strategy
	ShouldFallback(err error) bool

	// Configure configures the strategy with the given configuration
	Configure(cfg core.Configuration) error
}

type StrategyConfig struct {
	Timeout     int
	MaxRetries  int
	EnableProxy bool
	ProxyType   string // "http" or "socks5"
	ProxyAddr   string
	RewriteRule string
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
