package utils

import (
	"context"
	"runtime"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

//
var (
	ErrGithubReleaseAssetNotFound = errors.New("github release asset not found")
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock GithubRepositoryClient

type GithubRepositoryClient interface {
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
}

func GetCmdrRelease(ctx context.Context, client GithubRepositoryClient, releaseName string) (release *github.RepositoryRelease, err error) {
	if releaseName == "latest" {
		release, _, err = client.GetLatestRelease(ctx, core.Author, core.Name)
	} else {
		release, _, err = client.GetReleaseByTag(ctx, core.Author, core.Name, releaseName)
	}

	return
}

func SearchReleaseAsset(ctx context.Context, assetName string, release *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	assets := NewSortedHeap(len(release.Assets))
	for _, asset := range release.Assets {
		if asset.BrowserDownloadURL == nil {
			continue
		}

		currentAssetName := asset.GetName()

		if currentAssetName == assetName {
			assets.Add(asset, 0.0)
			break
		}

		score := 0.0
		if strings.Contains(currentAssetName, runtime.GOOS) {
			score += 1
		}
		if strings.Contains(currentAssetName, runtime.GOARCH) {
			score += 1
		}
		if score > 0.0 {
			assets.Add(asset, score)
		}
	}

	item, _ := assets.PopMax()
	if item == nil {
		return nil, errors.Wrapf(ErrGithubReleaseAssetNotFound, "search release asset failed")
	}

	return item.(*github.ReleaseAsset), nil
}
