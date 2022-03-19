package utils_test

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v39/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
	"github.com/mrlyc/cmdr/core/utils/mock"
)

var _ = Describe("Github", func() {
	var (
		ctrl     *gomock.Controller
		client   *mock.MockGithubRepositoryClient
		ctx      context.Context
		release  *github.RepositoryRelease
		searcher *utils.CmdrApiSearcher
	)

	BeforeEach(func() {
		release = &github.RepositoryRelease{}
		ctrl = gomock.NewController(GinkgoT())
		client = mock.NewMockGithubRepositoryClient(ctrl)
		ctx = context.Background()
		searcher = utils.NewCmdrApiSearcher(client)
	})

	AfterEach(func() {
		ctrl.Finish()
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

	It("GetLatestAsset", func() {
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

		asset, err := searcher.GetLatestAsset(ctx, releaseName, assetName)
		Expect(err).To(BeNil())

		Expect(asset.Name).To(Equal(releaseName))
		Expect(asset.Version).To(Equal("1.0.0"))
		Expect(asset.Asset).To(Equal(assetName))
		Expect(asset.Url).To(Equal(fakeUrl))
	})
})
