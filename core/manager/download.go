package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type DownloadManager struct {
	core.CommandManager
	fetcher core.Fetcher
}

func (m *DownloadManager) search(name, output string) (string, error) {
	type searchedFile struct {
		path  string
		score float64
	}

	files := make([]searchedFile, 0, 1)
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
			files = append(files, searchedFile{
				path:  path,
				score: score,
			})
		}

		return nil
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to walk %s", output)
	}

	if len(files) == 0 {
		return "", errors.Wrapf(core.ErrBinaryNotFound, "binary %s not found", name)
	}

	sort.SliceStable(files, func(i, j int) bool {
		if files[i].score != files[j].score {
			return files[i].score > files[j].score
		}

		return len(files[i].path) < len(files[j].path)
	})

	return files[0].path, nil
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
