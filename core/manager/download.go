package manager

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/fetcher"
	"github.com/mrlyc/cmdr/core/strategy"
	"github.com/mrlyc/cmdr/core/utils"
)

type DownloadManager struct {
	core.CommandManager
	fetchers     []core.Fetcher
	retries      int
	replacements utils.Replacements
	strategy     *strategy.StrategyChain
}

func (m *DownloadManager) SetReplacements(replacements utils.Replacements) {
	m.replacements = replacements
}

func (m *DownloadManager) SetStrategyChain(chain *strategy.StrategyChain) {
	m.strategy = chain
}

func (m *DownloadManager) getFetcherOptions() []getter.ClientOption {
	var options []getter.ClientOption

	// Get options from strategy chain
	if m.strategy != nil {
		// For now, use direct strategy as default
		// In future, we could select strategy based on URI
		directStrategy := strategy.NewDirectStrategy()
		directStrat, ok := directStrategy.(*strategy.DirectStrategy)
		if ok {
			options = append(options, directStrat.GetOptions()...)
		}
	}

	// Add default options if not set
	if len(options) == 0 {
		options = []getter.ClientOption{
			getter.WithTimeout(30 * time.Second),
		}
	}

	return options
}

func (m *DownloadManager) search(name, output string) (string, error) {
	files := utils.NewSortedHeap(1)
	nameLower := strings.ToLower(name)
	nameLength := float64(len(nameLower))

	err := filepath.Walk(output, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		score := 0.0
		if info.Mode()&0111 != 0 {
			score = 0.1 / nameLength // prefer to choose executable file
		}

		file := filepath.Base(path)
		if strings.Contains(strings.ToLower(file), nameLower) {
			score += nameLength / float64(len(file))
		}

		if score > 0 {
			files.Add(path, score)
		}

		return nil
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to walk %s", output)
	}

	if files.Len() == 0 {
		return "", errors.Wrapf(core.ErrBinaryNotFound, "binary %s not found", name)
	}

	file, _ := files.PopMax()

	return file.(string), nil
}

func (m *DownloadManager) fetch(fetcher core.Fetcher, name, version, location, output string) (string, error) {
	logger := core.GetLogger()
	logger.Info("fetching", map[string]interface{}{
		"uri": location,
	})

	// Use strategy chain if available
	if m.strategy != nil {
		var finalResult string
		var finalErr error

		// Execute strategy chain
		err := m.strategy.Execute(location, func(uri string) error {
			logger.Debug("downloading with URI", map[string]interface{}{
				"uri": uri,
			})

			// Apply URL rewriting if enabled
			for _, strat := range m.strategy.(*strategy.StrategyChain).Strategies() {
				if rewriteStrat, ok := strat.(*strategy.RewriteStrategy); ok && rewriteStrat.IsEnabled() {
					rewritten, err := rewriteStrat.GetRewrittenURI(uri)
					if err != nil {
						logger.Warn("URL rewrite failed, using original", map[string]interface{}{
							"error": err.Error(),
						})
					} else if rewritten != uri {
						logger.Info("URL rewritten", map[string]interface{}{
							"original": uri,
							"rewritten": rewritten,
						})
						uri = rewritten
					}
				}
			}

			// Apply replacements
			uri, _ = m.replacements.ReplaceString(uri)

			// Update fetcher options based on current strategy
			if gg, ok := fetcher.(*fetcher.GoGetter); ok {
				options := m.getFetcherOptions()
				gg.SetOptions(options)
			}

			// Try download
			fetchErr := fetcher.Fetch(name, version, uri, output)
			if fetchErr != nil {
				return fetchErr
			}

			// Download succeeded, search for binary
			result, searchErr := m.search(name, output)
			if searchErr != nil {
				return searchErr
			}

			finalResult = result
			return nil
		})

		if err != nil {
			return "", errors.Wrapf(err, "failed to download %s", location)
		}

		return finalResult, nil
	}

	// Fallback to old retry logic
	var err error
	for i := 0; i < m.retries; i++ {
		// Apply replacements
		location, _ = m.replacements.ReplaceString(location)

		err = fetcher.Fetch(name, version, location, output)
		if err == nil {
			break
		} else {
			logger.Warn("download failed, retrying...", map[string]interface{}{
				"uri": location,
			})
		}
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", location)
	}

	return m.search(name, output)
}
				}
			}

			// Apply replacements
			uri, _ = m.replacements.ReplaceString(uri)

			// Try download
			fetchErr := fetcher.Fetch(name, version, uri, output)
			if fetchErr != nil {
				return fetchErr
			}

			// Download succeeded, search for binary
			result, searchErr := m.search(name, output)
			if searchErr != nil {
				return searchErr
			}

			finalResult = result
			return nil
		})

		if err != nil {
			return "", errors.Wrapf(err, "failed to download %s", location)
		}

		return finalResult, nil
	}

	// Fallback to old retry logic
	var err error
	for i := 0; i < m.retries; i++ {
		// Apply replacements
		location, _ = m.replacements.ReplaceString(location)

		err = fetcher.Fetch(name, version, location, output)
		if err == nil {
			break
		} else {
			logger.Warn("download failed, retrying...", map[string]interface{}{
				"uri": location,
			})
		}
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", location)
	}

	return m.search(name, output)
}

func (m *DownloadManager) Define(name string, version string, uriOrLocation string) (core.Command, error) {
	uriOrLocation, _ = m.replacements.ReplaceString(uriOrLocation)

	for _, fetcher := range m.fetchers {
		if !fetcher.IsSupport(uriOrLocation) {
			continue
		}

		dst, err := os.MkdirTemp("", "")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create temp dir")
		}
		defer os.RemoveAll(dst)

		location, err := m.fetch(fetcher, name, version, uriOrLocation, dst)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to fetch %s", location)
		}

		uriOrLocation = location
	}

	return m.CommandManager.Define(name, version, uriOrLocation)
}

func NewDownloadManager(
	manager core.CommandManager, fetchers []core.Fetcher, retries int, replacements utils.Replacements,
) *DownloadManager {
	return &DownloadManager{
		CommandManager: manager,
		fetchers:       fetchers,
		retries:        retries,
		replacements:   replacements,
	}
}

func init() {
	core.RegisterCommandManagerFactory(core.CommandProviderDownload, func(cfg core.Configuration) (core.CommandManager, error) {
		manager, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			utils.ExitOnError("Failed to create command manager", err)
		}

		var replacements utils.Replacements
		err = cfg.UnmarshalKey(core.CfgKeyDownloadReplace, &replacements)
		if err != nil {
			utils.ExitOnError("Failed to parse download replace config", err)
		}

		// Create strategy chain
		strategyChain := strategy.NewStrategyChain(
			strategy.NewDirectStrategy(),
			strategy.NewRewriteStrategy(),
			strategy.NewProxyStrategy(),
		)

		// Configure strategies
		if err := strategyChain.Configure(cfg); err != nil {
			utils.ExitOnError("Failed to configure download strategies", err)
		}

		downloadManager := NewDownloadManager(manager, []core.Fetcher{
			fetcher.NewDefaultGoInstaller(),
			fetcher.NewDefaultGoGetter(os.Stderr),
		}, 3, replacements)

		downloadManager.SetStrategyChain(strategyChain)

		return downloadManager, nil
	})
}

		var replacements utils.Replacements

		err = cfg.UnmarshalKey(core.CfgKeyDownloadReplace, &replacements)
		if err != nil {
			utils.ExitOnError("Failed to parse download replace config", err)
		}

		return NewDownloadManager(manager, []core.Fetcher{
			fetcher.NewDefaultGoInstaller(),
			fetcher.NewDefaultGoGetter(os.Stderr),
		}, 3, replacements), nil
	})
}
