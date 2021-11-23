package core

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type Downloader struct {
	BaseStep
	schemaRegexp *regexp.Regexp
	tempDir      string
}

func (d *Downloader) String() string {
	return "downloader"
}

func (d *Downloader) Run(ctx context.Context) (context.Context, error) {
	name := utils.GetStringFromContext(ctx, define.ContextKeyName)
	url := utils.GetStringFromContext(ctx, define.ContextKeyLocation)
	logger := define.Logger
	var err error

	if !d.schemaRegexp.MatchString(url) {
		return ctx, nil
	}

	d.tempDir, err = afero.TempDir(define.FS, "", "")
	utils.ExitWithError(err, "create temporary dir failed")

	location := filepath.Join(d.tempDir, name)

	logger.Info("downloading", map[string]interface{}{
		"url":  url,
		"name": name,
	})
	err = utils.DownloadToFile(ctx, url, location)
	if err != nil {
		return ctx, err
	}

	logger.Info("command downloaded", map[string]interface{}{
		"url": url,
	})

	return context.WithValue(ctx, define.ContextKeyLocation, location), nil
}

func (d *Downloader) Finish(ctx context.Context) error {
	if d.tempDir == "" {
		return nil
	}

	return define.FS.RemoveAll(d.tempDir)
}

func NewDownloader() *Downloader {
	return &Downloader{
		schemaRegexp: regexp.MustCompile(`^https?://`),
	}
}
