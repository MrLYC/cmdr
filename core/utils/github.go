package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

var (
	ErrReleaseAssetNotFound = fmt.Errorf("release asset not found")
)

func GetCMDRRelease(ctx context.Context, tag string) (release *github.RepositoryRelease, err error) {
	client := github.NewClient(nil)
	if tag == "latest" {
		release, _, err = client.Repositories.GetLatestRelease(ctx, core.Author, core.Name)
	} else {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, core.Author, core.Name, tag)
	}

	return
}

func DownloadReleaseAssetByName(ctx context.Context, release *github.RepositoryRelease, name, output string) error {
	for _, asset := range release.Assets {
		if name != *asset.Name {
			continue
		}

		return DownloadToFile(ctx, *asset.BrowserDownloadURL, output)
	}

	return errors.Wrapf(ErrReleaseAssetNotFound, name)
}
