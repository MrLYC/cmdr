package strategy

import (
	"errors"
	"fmt"

	"github.com/mrlyc/cmdr/core"
)

type StrategyChain struct {
	strategies []DownloadStrategy
	config     *StrategyConfig
}

func (c *StrategyChain) Strategies() []DownloadStrategy {
	return c.strategies
}

func (c *StrategyChain) AddStrategy(strategy DownloadStrategy) {
	c.strategies = append(c.strategies, strategy)
}

func (c *StrategyChain) GetEnabledStrategies(uri string) []DownloadStrategy {
	logger := core.GetLogger()
	var enabled []DownloadStrategy

	for _, strategy := range c.strategies {
		if strategy.IsEnabled(uri) {
			logger.Debug("strategy enabled for URI", map[string]interface{}{
				"strategy": strategy.Name(),
				"uri":      uri,
			})
			enabled = append(enabled, strategy)
		}
	}

	return enabled
}

func (c *StrategyChain) Execute(uri string, downloadFunc func(string) error) error {
	logger := core.GetLogger()
	var lastErr error

	// Get strategies that are enabled for this URI
	enabledStrategies := c.GetEnabledStrategies(uri)
	if len(enabledStrategies) == 0 {
		// No strategy enabled, try all strategies (backward compatibility)
		enabledStrategies = c.strategies
	}

	for strategyIdx, strategy := range enabledStrategies {
		strategyName := strategy.Name()

		// Prepare URI with strategy
		preparedURI, err := strategy.Prepare(uri)
		if err != nil {
			logger.Warn("strategy prepare failed, trying next", map[string]interface{}{
				"strategy": strategyName,
				"error":    err.Error(),
			})
			continue
		}

		logger.Info("using download strategy", map[string]interface{}{
			"strategy": strategyName,
			"uri":      preparedURI,
		})

		// Try download with this strategy
		retryCount := 0
		maxRetries := c.getStrategyMaxRetries(strategy)

		for retryCount < maxRetries {
			err = downloadFunc(preparedURI)

			if err == nil {
				logger.Info("download succeeded", map[string]interface{}{
					"strategy": strategyName,
					"retries":  retryCount,
				})
				return nil
			}

			retryCount++
			lastErr = err

			// Check if we should retry with same strategy
			if retryCount < maxRetries && strategy.ShouldRetry(err) {
				logger.Warn("download failed, retrying with same strategy", map[string]interface{}{
					"strategy": strategyName,
					"retry":    retryCount,
					"error":    err.Error(),
				})
				continue
			}

			// Check if we should fallback to next strategy
			if strategy.ShouldFallback(err) {
				logger.Warn("download failed, trying next strategy", map[string]interface{}{
					"strategy": strategyName,
					"error":    err.Error(),
				})
				break
			}

			// Error not retriable or fallback-able
			logger.Error("download failed with non-retriable error", map[string]interface{}{
				"strategy": strategyName,
				"error":    err.Error(),
			})
			return err
		}

		// If this is not the last strategy and we have error, continue to next
		if strategyIdx < len(enabledStrategies)-1 {
			nextStrategy := enabledStrategies[strategyIdx+1]
			logger.Info("switching to next strategy", map[string]interface{}{
				"current": strategyName,
				"next":    nextStrategy.Name(),
			})
			continue
		}

		// Last strategy failed
		break
	}

	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrAllStrategiesFailed, lastErr)
	}

	return errors.New("unexpected state: no error but download failed")
}

func (c *StrategyChain) getStrategyMaxRetries(strategy DownloadStrategy) int {
	// Default to 3 retries if not configured
	if c.config != nil && c.config.MaxRetries > 0 {
		return c.config.MaxRetries
	}
	return 3
}

func (c *StrategyChain) Configure(cfg core.Configuration) error {
	for _, strategy := range c.strategies {
		if err := strategy.Configure(cfg); err != nil {
			return fmt.Errorf("failed to configure strategy %s: %w", strategy.Name(), err)
		}
	}
	return nil
}

func NewStrategyChain(strategies ...DownloadStrategy) *StrategyChain {
	return &StrategyChain{
		strategies: strategies,
	}
}
