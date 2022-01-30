package operator

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type Downloader struct {
	BaseOperator
	schemaRegexp *regexp.Regexp
	tempDir      string
}

func (d *Downloader) String() string {
	return "downloader"
}

func (d *Downloader) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	var errs error
	for _, command := range commands {
		name := command.Name
		url := command.Location
		if !d.schemaRegexp.MatchString(url) {
			continue
		}

		d.tempDir, err = afero.TempDir(define.FS, "", "")
		utils.ExitWithError(err, "create temporary dir failed")

		location := filepath.Join(d.tempDir, name)
		command.Location = location

		logger.Info("downloading", map[string]interface{}{
			"url":  url,
			"name": name,
		})
		err = utils.DownloadToFile(ctx, url, location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "download %s failed", url))
			continue
		}

		logger.Info("command downloaded", map[string]interface{}{
			"url": url,
		})
	}

	return ctx, errs
}

func (d *Downloader) Commit(ctx context.Context) error {
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
