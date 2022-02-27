package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type DownloadManager struct {
	core.CommandManager
	fetcher core.Fetcher
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

func (m *DownloadManager) fetch(name, version, location, output string) (string, error) {
	err := m.fetcher.Fetch(location, output)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", location)
	}

	return m.search(name, output)
}

func (m *DownloadManager) Define(name string, version string, uri string) error {
	if !m.fetcher.IsSupport(uri) {
		return m.CommandManager.Define(name, version, uri)
	}

	dst, err := os.MkdirTemp("", "")
	if err != nil {
		return errors.Wrapf(err, "failed to create temp dir")
	}
	defer os.RemoveAll(dst)

	location, err := m.fetch(name, version, uri, dst)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch %s", location)
	}

	return m.CommandManager.Define(name, version, location)
}

func NewDownloadManager(manager core.CommandManager, fetcher core.Fetcher) *DownloadManager {
	return &DownloadManager{
		CommandManager: manager,
		fetcher:        fetcher,
	}
}

func init() {
	core.RegisterCommandManagerFactory(core.CommandProviderDownload, func(cfg core.Configuration) (core.CommandManager, error) {
		manager, err := core.NewCommandManager(core.CommandProviderDefault, cfg)
		if err != nil {
			utils.ExitOnError("Failed to create command manager", err)
		}

		return NewDownloadManager(manager, utils.NewProgressBarDownloader(os.Stderr)), nil
	})
}
