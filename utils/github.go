package utils

import (
	"context"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
)

var (
	ErrReleaseAssertNotFound = errors.New("release assert not found")
)

func GetCMDRRelease(ctx context.Context, tag string) (release *github.RepositoryRelease, err error) {
	client := github.NewClient(nil)
	if tag == "latest" {
		release, _, err = client.Repositories.GetLatestRelease(ctx, define.Author, define.Name)
	} else {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, define.Author, define.Name, tag)
	}

	return
}

func DownloadReleaseAssertByName(ctx context.Context, release *github.RepositoryRelease, name, output string) error {
	for _, assert := range release.Assets {
		if name != *assert.Name {
			continue
		}

		return DownloadToFile(ctx, *assert.BrowserDownloadURL, output)
	}

	return errors.Wrapf(ErrReleaseAssertNotFound, name)
}
