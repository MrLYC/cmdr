package operator

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type ReleaseSearcher struct {
	BaseOperator
	release string
	asset   string
}

func (r *ReleaseSearcher) String() string {
	return "release-searcher"
}

func (r *ReleaseSearcher) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	command, err := GetCommandFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get command from context failed")
	}

	logger.Info("searching release", map[string]interface{}{
		"release": r.release,
		"asset":   r.asset,
		"name":    command.Name,
	})
	release, err := utils.GetCMDRRelease(ctx, r.release)
	if err != nil {
		return ctx, err
	}

	version := strings.TrimPrefix(*release.TagName, "v")
	logger.Debug("searching asset", map[string]interface{}{
		"version": version,
		"name":    command.Name,
	})

	var url string
	for _, asset := range release.Assets {
		if r.asset != *asset.Name {
			logger.Debug("skip asset", map[string]interface{}{
				"name":     *asset.Name,
				"excepted": r.asset,
			})
			continue
		}

		url = *asset.BrowserDownloadURL
	}

	if url == "" {
		return ctx, errors.Wrapf(ErrAssetNotFound, "release not found: %s", r.asset)
	}

	command.Version = version
	command.Location = url

	logger.Info("release found", map[string]interface{}{
		"name":    command.Name,
		"release": version,
		"url":     url,
	})
	return ctx, nil
}

func NewReleaseSearcher(release, asset string) *ReleaseSearcher {
	return &ReleaseSearcher{
		release: release,
		asset:   asset,
	}
}
