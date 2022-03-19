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

type CmdrApiSearcher struct {
	client GithubRepositoryClient
}

func (s *CmdrApiSearcher) SearchReleaseAsset(ctx context.Context, assetName string, release *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	assets := NewSortedHeap(len(release.Assets))
	for _, asset := range release.Assets {
		if asset.BrowserDownloadURL == nil {
			continue
		}

		currentAssetName := asset.GetName()

		if currentAssetName == assetName {
			return asset, nil
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

func (s *CmdrApiSearcher) GetCmdrRelease(ctx context.Context, releaseName string) (release *github.RepositoryRelease, err error) {
	if releaseName == "latest" {
		release, _, err = s.client.GetLatestRelease(ctx, core.Author, core.Name)
	} else {
		release, _, err = s.client.GetReleaseByTag(ctx, core.Author, core.Name, releaseName)
	}

	return
}

func (s *CmdrApiSearcher) GetLatestAsset(ctx context.Context, releaseName, assetName string) (result core.CmdrReleaseInfo, err error) {
	logger := core.GetLogger()

	logger.Info("searching for release", map[string]interface{}{
		"release": releaseName,
	})
	release, err := s.GetCmdrRelease(ctx, releaseName)
	if err != nil {
		return result, errors.Wrapf(err, "search release failed")
	}

	asset, err := s.SearchReleaseAsset(ctx, assetName, release)
	if err != nil {
		return result, errors.Wrapf(err, "search release asset failed")
	}

	result.Name = release.GetName()
	result.Version = strings.Trim(release.GetTagName(), "v")
	result.Asset = asset.GetName()
	result.Url = asset.GetBrowserDownloadURL()

	logger.Info("release asset found", map[string]interface{}{
		"release": result.Name,
		"asset":   result.Asset,
	})

	return result, nil
}

func NewCmdrApiSearcher(client GithubRepositoryClient) *CmdrApiSearcher {
	return &CmdrApiSearcher{
		client: client,
	}
}

func init() {
	core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderApi, func(cfg core.Configuration) (core.CmdrSearcher, error) {
		return NewCmdrApiSearcher(github.NewClient(nil).Repositories), nil
	})
}
