package utils_test

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v39/github"
	"github.com/mmcdole/gofeed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	coremock "github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/core/utils"
	"github.com/mrlyc/cmdr/core/utils/mock"
)

var _ = Describe("Github", func() {
	var (
		ctx  context.Context
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CmdrApiFetcher", func() {
		var (
			client   *mock.MockGithubRepositoryClient
			release  *github.RepositoryRelease
			searcher *utils.CmdrApiFetcher
		)

		BeforeEach(func() {
			release = &github.RepositoryRelease{}
			client = mock.NewMockGithubRepositoryClient(ctrl)
			searcher = utils.NewCmdrApiFetcher(client)
		})

		It("should get latest release", func() {
			client.EXPECT().GetLatestRelease(ctx, core.Author, core.Name).Return(release, nil, nil)

			_, err := searcher.GetCmdrRelease(ctx, "latest")
			Expect(err).To(BeNil())
		})

		It("should get named release", func() {
			client.EXPECT().GetReleaseByTag(ctx, core.Author, core.Name, "v0.0.0").Return(release, nil, nil)

			_, err := searcher.GetCmdrRelease(ctx, "v0.0.0")
			Expect(err).To(BeNil())
		})

		Context("SearchReleaseAsset", func() {
			var (
				assetName = "cmdr-goos-goarch"
			)

			DescribeTable("should detect asset", func(expected string, assetPatterns [][]string) {
				fakeUrl := "http://example.com"
				assets := make([]*github.ReleaseAsset, 0, len(assetPatterns))

				for _, p := range assetPatterns {
					name := strings.Join(p, "-")
					assets = append(assets, &github.ReleaseAsset{
						BrowserDownloadURL: &fakeUrl,
						Name:               &name,
					})
				}

				release := &github.RepositoryRelease{
					Assets: assets,
				}

				asset, err := searcher.SearchReleaseAsset(context.Background(), assetName, release)
				Expect(err).To(BeNil())
				Expect(asset.GetName()).To(Equal(expected))
			},
				Entry("same asset name", assetName, [][]string{
					{assetName},
					{"foo", "bar"},
				}),
				Entry("by os", runtime.GOOS, [][]string{
					{"foo", "bar"},
					{runtime.GOOS},
				}),
				Entry("by arch", runtime.GOARCH, [][]string{
					{"foo", "bar"},
					{runtime.GOARCH},
				}),
				Entry("by os and arch", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH), [][]string{
					{"foo", "bar"},
					{runtime.GOOS, "baz"},
					{runtime.GOARCH, "qux"},
					{runtime.GOOS, runtime.GOARCH},
				}),
				Entry("prefer to asset name", assetName, [][]string{
					{assetName},
					{runtime.GOOS, runtime.GOARCH},
				}),
			)
		})

		It("GetReleaseAsset", func() {
			releaseName := "latest"
			assetName := "cmdr-goos-goarch"
			fakeUrl := "http://example.com"
			tagName := "v1.0.0"
			release.Name = &releaseName
			release.TagName = &tagName
			release.Assets = []*github.ReleaseAsset{
				{
					Name:               &assetName,
					BrowserDownloadURL: &fakeUrl,
				},
			}
			client.EXPECT().GetLatestRelease(ctx, core.Author, core.Name).Return(release, nil, nil)

			asset, err := searcher.GetReleaseAsset(ctx, releaseName, assetName)
			Expect(err).To(BeNil())

			Expect(asset.Name).To(Equal(releaseName))
			Expect(asset.Version).To(Equal("1.0.0"))
			Expect(asset.Asset).To(Equal(assetName))
			Expect(asset.Url).To(Equal(fakeUrl))
		})
	})

	Context("CmdrFeedFetcher", func() {
		var (
			feed     gofeed.Feed
			searcher *utils.CmdrFeedFetcher
		)

		BeforeEach(func() {
			now := time.Now()
			future := now.Add(time.Hour)

			feed = gofeed.Feed{
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.1",
						PublishedParsed: &future,
					},
					{
						Title:           "v1.0.0",
						PublishedParsed: &now,
					},
				},
			}
			searcher = utils.NewCmdrFeedFetcher(func(ctx context.Context) (*gofeed.Feed, error) {
				return &feed, nil
			})
		})

		It("should return latest release", func() {
			info, err := searcher.GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch")
			Expect(err).To(BeNil())

			Expect(info.Name).To(Equal("v1.0.1"))
			Expect(info.Version).To(Equal("1.0.1"))
			Expect(info.Asset).To(Equal("cmdr-goos-goarch"))
			Expect(info.Url).To(Equal(
				"https://github.com/MrLYC/cmdr/releases/download/v1.0.1/cmdr-goos-goarch",
			))
		})

		It("should return specified release", func() {
			info, err := searcher.GetReleaseAsset(ctx, "v1.0.0", "cmdr-goos-goarch")
			Expect(err).To(BeNil())

			Expect(info.Name).To(Equal("v1.0.0"))
			Expect(info.Version).To(Equal("1.0.0"))
			Expect(info.Asset).To(Equal("cmdr-goos-goarch"))
			Expect(info.Url).To(Equal(
				"https://github.com/MrLYC/cmdr/releases/download/v1.0.0/cmdr-goos-goarch",
			))
		})
	})

	Context("CmdrReleaseSearcher", func() {
		var (
			searcher                     *utils.CmdrReleaseSearcher
			mockSearcher1, mockSearcher2 *coremock.MockCmdrSearcher
		)

		BeforeEach(func() {
			mockSearcher1 = coremock.NewMockCmdrSearcher(ctrl)
			mockSearcher2 = coremock.NewMockCmdrSearcher(ctrl)
			searcher = utils.NewCmdrReleaseSearcher(mockSearcher1, mockSearcher2)
		})

		It("should return latest release from searcher1", func() {
			var release core.CmdrReleaseAsset
			mockSearcher1.
				EXPECT().
				GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch").
				Return(release, nil)

			info, err := searcher.GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch")
			Expect(err).To(BeNil())
			Expect(info).To(Equal(release))
		})

		It("should return latest release from searcher2", func() {
			var release1, release2 core.CmdrReleaseAsset

			mockSearcher1.
				EXPECT().
				GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch").
				Return(release1, fmt.Errorf("testing"))

			mockSearcher2.
				EXPECT().
				GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch").
				Return(release2, nil)

			info, err := searcher.GetReleaseAsset(ctx, "latest", "cmdr-goos-goarch")
			Expect(err).To(BeNil())
			Expect(info).To(Equal(release2))
		})
	})
})
