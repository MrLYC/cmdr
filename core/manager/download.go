package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/fetcher"
	"github.com/mrlyc/cmdr/core/utils"
)

type DownloadManager struct {
	core.CommandManager
	fetchers []core.Fetcher
	retries  int
}

func (m *DownloadManager) search(name, output string) (string, error) {
	files := utils.NewSortedHeap(1)
	nameLength := float64(len(name))

	err := filepath.Walk(output, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		score := 0.0
		if info.Mode()&0111 != 0 {
			score = 0.1 / nameLength // perfer to choose executable file
		}

		file := filepath.Base(path)
		if strings.Contains(file, name) {
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
	var err error
	logger := core.GetLogger()

	for i := 0; i < m.retries; i++ {
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

func NewDownloadManager(manager core.CommandManager, fetchers []core.Fetcher, retries int) *DownloadManager {
	return &DownloadManager{
		CommandManager: manager,
		fetchers:       fetchers,
		retries:        retries,
	}
}

func init() {
	core.RegisterCommandManagerFactory(core.CommandProviderDownload, func(cfg core.Configuration) (core.CommandManager, error) {
		manager, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			utils.ExitOnError("Failed to create command manager", err)
		}

		return NewDownloadManager(manager, []core.Fetcher{
			fetcher.NewGoInstaller(),
			fetcher.NewDefaultGoGetter(os.Stderr),
		}, 3), nil
	})
}
