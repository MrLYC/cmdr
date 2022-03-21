package utils

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
	"github.com/mmcdole/gofeed"
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

type CmdrApiFetcher struct {
	client GithubRepositoryClient
}

func (s *CmdrApiFetcher) String() string {
	return "github-api"
}

func (s *CmdrApiFetcher) SearchReleaseAsset(ctx context.Context, assetName string, release *github.RepositoryRelease) (*github.ReleaseAsset, error) {
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

func (s *CmdrApiFetcher) GetCmdrRelease(ctx context.Context, releaseName string) (release *github.RepositoryRelease, err error) {
	if releaseName == "latest" {
		release, _, err = s.client.GetLatestRelease(ctx, core.Author, core.Name)
	} else {
		release, _, err = s.client.GetReleaseByTag(ctx, core.Author, core.Name, releaseName)
	}

	return
}

func (s *CmdrApiFetcher) GetReleaseAsset(ctx context.Context, releaseName, assetName string) (result core.CmdrReleaseAsset, err error) {
	logger := core.GetLogger()
	logger.Debug("searching cmdr release by github api", map[string]interface{}{
		"release": releaseName,
		"asset":   assetName,
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

func NewCmdrApiFetcher(client GithubRepositoryClient) *CmdrApiFetcher {
	return &CmdrApiFetcher{
		client: client,
	}
}

type CmdrFeedFetcher struct {
	fetchFn func(ctx context.Context) (feed *gofeed.Feed, err error)
}

func (s *CmdrFeedFetcher) String() string {
	return "github-feed"
}

func (s *CmdrFeedFetcher) searchRelease(releaseName string, feed *gofeed.Feed) *gofeed.Item {
	if releaseName == "latest" {
		return feed.Items[len(feed.Items)-1]
	}

	for _, item := range feed.Items {
		if item.Title == releaseName {
			return item
		}
	}

	return nil
}

func (s *CmdrFeedFetcher) getRelease(ctx context.Context, releaseName string) (*gofeed.Item, error) {
	feed, err := s.fetchFn(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch cmdr atom feed failed")
	}

	sort.Sort(feed)
	item := s.searchRelease(releaseName, feed)
	if item == nil {
		return nil, errors.Wrapf(ErrGithubReleaseAssetNotFound, "search release %s failed", releaseName)
	}

	return item, nil
}

func (s *CmdrFeedFetcher) GetRelease(ctx context.Context, releaseName string) (string, error) {
	item, err := s.getRelease(ctx, releaseName)
	if err != nil {
		return "", err
	}

	return item.Title, nil
}

func (s *CmdrFeedFetcher) GetReleaseAsset(ctx context.Context, releaseName, assetName string) (result core.CmdrReleaseAsset, err error) {
	logger := core.GetLogger()
	logger.Debug("searching cmdr release by github feed", map[string]interface{}{
		"release": releaseName,
		"asset":   assetName,
	})

	item, err := s.getRelease(ctx, releaseName)
	if err != nil {
		return result, err
	}

	releaseVersion, err := version.NewVersion(item.Title)
	if err != nil {
		return result, errors.Wrapf(err, "parse release %s version failed", item.Title)
	}

	logger.Info("release asset found", map[string]interface{}{
		"release":    item.Title,
		"url":        item.Link,
		"publish_at": item.Published,
		"update_at":  item.Updated,
	})

	result.Name = item.Title
	result.Version = releaseVersion.String()
	result.Asset = assetName
	result.Url = fmt.Sprintf(
		`https://github.com/MrLYC/cmdr/releases/download/%s/%s`,
		item.Title, assetName,
	)

	logger.Info("release found", map[string]interface{}{
		"release": result.Name,
	})

	return
}

func NewCmdrFeedFetcher(fetchFn func(ctx context.Context) (feed *gofeed.Feed, err error)) *CmdrFeedFetcher {
	return &CmdrFeedFetcher{
		fetchFn: fetchFn,
	}
}

func NewCmdrAtomFetcher() *CmdrFeedFetcher {
	return NewCmdrFeedFetcher(func(ctx context.Context) (feed *gofeed.Feed, err error) {
		return gofeed.NewParser().ParseURL(
			fmt.Sprintf(
				`https://github.com/%s/%s/releases.atom`,
				core.Author, core.Name,
			),
		)
	})
}

type CmdrReleaseSearcher struct {
	searchers []core.CmdrSearcher
}

func (s *CmdrReleaseSearcher) GetReleaseAsset(ctx context.Context, releaseName, assetName string) (result core.CmdrReleaseAsset, err error) {
	logger := core.GetLogger()
	var errs error

	for _, searcher := range s.searchers {
		logger.Debug("searching cmdr release", map[string]interface{}{
			"searcher": searcher,
			"release":  releaseName,
			"asset":    assetName,
		})

		result, err = searcher.GetReleaseAsset(ctx, releaseName, assetName)
		if err == nil {
			return
		}

		logger.Debug("search cmdr release failed, continue", map[string]interface{}{
			"searcher": searcher,
			"error":    err,
		})
		errs = multierror.Append(errs, err)
	}

	return result, errs
}

func NewCmdrReleaseSearcher(searchers ...core.CmdrSearcher) *CmdrReleaseSearcher {
	return &CmdrReleaseSearcher{
		searchers: searchers,
	}
}

func init() {
	core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderApi, func(cfg core.Configuration) (core.CmdrSearcher, error) {
		return NewCmdrApiFetcher(github.NewClient(nil).Repositories), nil
	})

	core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderAtom, func(cfg core.Configuration) (core.CmdrSearcher, error) {
		return NewCmdrAtomFetcher(), nil
	})

	core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderDefault, func(cfg core.Configuration) (core.CmdrSearcher, error) {
		apiSearcher, err := core.NewCmdrSearcher(core.CmdrSearcherProviderApi, cfg)
		if err != nil {
			return nil, err
		}

		atomSearcher, err := core.NewCmdrSearcher(core.CmdrSearcherProviderAtom, cfg)
		if err != nil {
			return nil, err
		}

		return NewCmdrReleaseSearcher(apiSearcher, atomSearcher), nil
	})
}
