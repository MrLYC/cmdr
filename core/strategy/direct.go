package strategy

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/mrlyc/cmdr/core"
)

var (
	ErrNetworkError    = errors.New("network error")
	ErrTimeoutError    = errors.New("timeout error")
	ErrConnectionError = errors.New("connection error")
)

type DirectStrategy struct {
	config *StrategyConfig
}

func (s *DirectStrategy) Name() string {
	return "direct"
}

func (s *DirectStrategy) Prepare(uri string) (string, error) {
	return uri, nil
}

func (s *DirectStrategy) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	// Retry on timeout and connection errors
	if isTimeoutError(err) || isConnectionError(err) {
		return true
	}

	// Retry on temporary network errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Temporary() {
			return true
		}
	}

	return false
}

func (s *DirectStrategy) ShouldFallback(err error) bool {
	if err == nil {
		return false
	}

	// Fallback on network errors
	return isNetworkError(err)
}

func (s *DirectStrategy) Configure(cfg core.Configuration) error {
	s.config = &StrategyConfig{
		Timeout:    cfg.GetInt("download.direct.timeout"),
		MaxRetries: cfg.GetInt("download.direct.max_retries"),
	}

	if s.config.Timeout == 0 {
		s.config.Timeout = 30 // default 30 seconds
	}
	if s.config.MaxRetries == 0 {
		s.config.MaxRetries = 3 // default 3 retries
	}

	return nil
}

func (s *DirectStrategy) GetOptions() []getter.ClientOption {
	if s.config == nil {
		return nil
	}

	var options []getter.ClientOption

	// Set timeout
	if s.config.Timeout > 0 {
		timeout := time.Duration(s.config.Timeout) * time.Second
		options = append(options, getter.WithTimeout(timeout))
	}

	return options
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadline exceeded") ||
		strings.Contains(errMsg, "timed out")
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "connection reset") ||
		strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "network is unreachable") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "dial tcp")
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	return isTimeoutError(err) || isConnectionError(err)
}

func NewDirectStrategy() *DirectStrategy {
	return &DirectStrategy{}
}
