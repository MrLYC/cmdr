package operator

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type Downloader struct {
	BaseOperator
	tempDir string
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
		locationUrl := command.Location
		urlInfo, err := url.Parse(locationUrl)
		if err != nil || urlInfo.Scheme == "" {
			continue
		}

		d.tempDir, err = os.MkdirTemp("", "")
		utils.ExitWithError(err, "create temporary dir failed")

		location := filepath.Join(d.tempDir, name)
		command.Location = location

		logger.Info("downloading", map[string]interface{}{
			"url":  locationUrl,
			"name": name,
		})
		err = utils.DownloadToFile(ctx, locationUrl, location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "download %s failed", locationUrl))
			continue
		}

		logger.Info("command downloaded", map[string]interface{}{
			"url": locationUrl,
		})
	}

	return ctx, errs
}

func (d *Downloader) Commit(ctx context.Context) error {
	if d.tempDir == "" {
		return nil
	}

	return os.RemoveAll(d.tempDir)
}

func NewDownloader() *Downloader {
	return &Downloader{}
}
