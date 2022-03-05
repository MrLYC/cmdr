package utils_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v39/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
	"github.com/mrlyc/cmdr/core/utils/mock"
)

var _ = Describe("Github", func() {
	Context("GetRelease", func() {
		var (
			ctrl    *gomock.Controller
			client  *mock.MockGithubRepositoryClient
			ctx     context.Context
			release github.RepositoryRelease
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			client = mock.NewMockGithubRepositoryClient(ctrl)
			ctx = context.Background()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should get latest release", func() {
			client.EXPECT().GetLatestRelease(ctx, core.Author, core.Name).Return(&release, nil, nil)

			_, err := utils.GetCmdrRelease(ctx, client, "latest")
			Expect(err).To(BeNil())
		})

		It("should get named release", func() {
			client.EXPECT().GetReleaseByTag(ctx, core.Author, core.Name, "v0.0.0").Return(&release, nil, nil)

			_, err := utils.GetCmdrRelease(ctx, client, "v0.0.0")
			Expect(err).To(BeNil())
		})
	})
})
