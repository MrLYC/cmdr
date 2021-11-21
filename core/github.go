package core

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type ReleaseSearcher struct {
	BaseStep
	release string
	asset   string
}

func (r *ReleaseSearcher) String() string {
	return "release-searcher"
}

func (r *ReleaseSearcher) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	logger.Debug("searching cmdr release", map[string]interface{}{
		"release": r.release,
		"asset":   r.asset,
	})
	release, err := utils.GetCMDRRelease(ctx, r.release)
	if err != nil {
		return ctx, err
	}

	version := strings.TrimPrefix(*release.TagName, "v")
	logger.Debug("searching cmdr asset", map[string]interface{}{
		"version": version,
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
		return ctx, errors.Wrapf(ErrAssetNotFound, "cmdr release not found: %s", r.asset)
	}

	return utils.SetIntoContext(ctx, map[define.ContextKey]interface{}{
		define.ContextKeyVersion:  version,
		define.ContextKeyLocation: url,
	}), nil
}

func NewReleaseSearcher(release, asset string) *ReleaseSearcher {
	return &ReleaseSearcher{
		release: release,
		asset:   asset,
	}
}
