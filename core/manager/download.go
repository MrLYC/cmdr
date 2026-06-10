package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func (m *DownloadManager) search(name, output string) (string, error) {
	files := utils.NewSortedHeap(1)
	nameLower := strings.ToLower(name)
	nameLength := float64(len(nameLower))

	err := filepath.WalkDir(output, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
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

func (m *DownloadManager) maybeRewriteURI(uri, original string) string {
	logger := core.GetLogger()
	if uri != original {
		return uri
	}

	for _, strat := range m.strategy.Strategies() {
		rewriteStrat, ok := strat.(*strategy.RewriteStrategy)
		if !ok || !rewriteStrat.IsEnabledConfigured() {
			continue
		}

		rewritten, err := rewriteStrat.GetRewrittenURI(uri)
		if err != nil {
			logger.Warn("URL rewrite failed, using original", map[string]interface{}{
				"error": err.Error(),
			})
			return uri
		}
		if rewritten != uri {
			logger.Info("URL rewritten", map[string]interface{}{
				"original":  uri,
				"rewritten": rewritten,
			})
			return rewritten
		}
	}

	return uri
}

func resetFetcherOptions(source core.Fetcher) {
	if gg, ok := source.(*fetcher.GoGetter); ok {
		gg.SetOptions(nil)
	}
}

func (m *DownloadManager) fetchOnce(f core.Fetcher, name, version, uri, output string) (string, error) {
	resetFetcherOptions(f)
	if err := f.Fetch(name, version, uri, output); err != nil {
		return "", err
	}

	return m.search(name, output)
}

func (m *DownloadManager) fetchWithStrategy(f core.Fetcher, name, version, location, output string) (string, error) {
	logger := core.GetLogger()
	var finalResult string

	err := m.strategy.Execute(location, func(uri string) error {
		logger.Debug("downloading with URI", map[string]interface{}{
			"uri": uri,
		})

		uri = m.maybeRewriteURI(uri, location)
		uri, _ = m.replacements.ReplaceString(uri)

		result, err := m.fetchOnce(f, name, version, uri, output)
		if err != nil {
			return err
		}

		finalResult = result
		return nil
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", location)
	}

	return finalResult, nil
}

func (m *DownloadManager) fetchWithRetries(f core.Fetcher, name, version, location, output string) (string, error) {
	logger := core.GetLogger()
	var err error
	uri := location

	for i := 0; i < m.retries; i++ {
		uri, _ = m.replacements.ReplaceString(uri)
		result, fetchErr := m.fetchOnce(f, name, version, uri, output)
		err = fetchErr
		if err == nil {
			return result, nil
		}

		logger.Warn("download failed, retrying...", map[string]interface{}{
			"uri": uri,
		})
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", location)
	}

	return "", errors.Wrapf(core.ErrBinaryNotFound, "binary %s not found", name)
}

func (m *DownloadManager) fetch(f core.Fetcher, name, version, location, output string) (string, error) {
	core.GetLogger().Info("fetching", map[string]interface{}{
		"uri": location,
	})

	if m.strategy != nil {
		return m.fetchWithStrategy(f, name, version, location, output)
	}

	return m.fetchWithRetries(f, name, version, location, output)
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
