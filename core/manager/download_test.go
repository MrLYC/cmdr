package manager_test

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/core/strategy"
	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Download", func() {
	var (
		ctrl      *gomock.Controller
		db        *mock.MockDatabase
		dbFactory func() (core.Database, error)
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDatabase(ctrl)
		core.SetDatabaseFactory(func() (core.Database, error) {
			return db, nil
		})
		dbFactory = core.GetDatabaseFactory()
	})

	AfterEach(func() {
		ctrl.Finish()
		core.SetDatabaseFactory(dbFactory)
	})

	Context("DownloadManager", func() {
		var (
			fetcher         *mock.MockFetcher
			baseManager     *mock.MockCommandManager
			downloadManager *manager.DownloadManager
			name            = "cmdr"
			version         = "1.0.0"
			uri             = ""
		)

		BeforeEach(func() {
			fetcher = mock.NewMockFetcher(ctrl)
			baseManager = mock.NewMockCommandManager(ctrl)
			downloadManager = manager.NewDownloadManager(baseManager, []core.Fetcher{fetcher}, 1, nil)
		})

		It("should call base manager", func() {
			fetcher.EXPECT().IsSupport(uri).Return(false)
			baseManager.EXPECT().Define(name, version, uri)

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		})

		It("should return base manager errors", func() {
			fetcher.EXPECT().IsSupport(uri).Return(false)
			baseManager.EXPECT().Define(name, version, uri).Return(nil, errors.New("define failed"))

			_, err := downloadManager.Define(name, version, uri)
			Expect(err).To(MatchError("define failed"))
		})

		It("should call with downloaded file", func() {
			var targetPath string

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, gomock.Any(), gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(os.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())

				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(targetPath).To(Equal(location))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		})

		It("should return fetch errors after retries", func() {
			downloadManager = manager.NewDownloadManager(baseManager, []core.Fetcher{fetcher}, 2, nil)
			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, uri, gomock.Any()).Return(errors.New("download failed")).Times(2)

			_, err := downloadManager.Define(name, version, uri)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to fetch"))
			Expect(err.Error()).To(ContainSubstring("failed to download"))
		})

		It("should return an error when downloaded output has no binary candidate", func() {
			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, uri, gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				Expect(os.MkdirAll(filepath.Join(dir, "empty-dir"), 0755)).To(Succeed())
				return nil
			})

			_, err := downloadManager.Define(name, version, uri)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("binary cmdr not found"))
		})

		It("should download through a strategy chain", func() {
			var targetPath string
			chain := strategy.NewStrategyChain(strategy.NewDirectStrategy())
			downloadManager.SetStrategyChain(chain)

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, uri, gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(os.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())
				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(location).To(Equal(targetPath))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		})

		It("should apply configured rewrite strategy while direct strategy downloads", func() {
			input := "https://github.com/MrLYC/cmdr"
			rewritten := "https://mirror.local/MrLYC/cmdr"
			var targetPath string

			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "https://mirror.local{{ .Path }}")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())
			chain := strategy.NewStrategyChain(strategy.NewDirectStrategy(), rewrite)
			downloadManager.SetStrategyChain(chain)

			fetcher.EXPECT().IsSupport(input).Return(true)
			fetcher.EXPECT().Fetch(name, version, rewritten, gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(os.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())
				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(location).To(Equal(targetPath))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, input)).To(Succeed())
		})

		It("should not rewrite twice when rewrite strategy prepares the URI", func() {
			input := "https://github.com/MrLYC/cmdr"
			rewritten := "https://mirror.local/MrLYC/cmdr"
			var targetPath string

			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "https://mirror.local{{ .Path }}")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())
			chain := strategy.NewStrategyChain(rewrite)
			downloadManager.SetStrategyChain(chain)

			fetcher.EXPECT().IsSupport(input).Return(true)
			fetcher.EXPECT().Fetch(name, version, rewritten, gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(os.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())
				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(location).To(Equal(targetPath))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, input)).To(Succeed())
		})

		It("should keep original URI when configured rewrite execution fails", func() {
			input := "https://github.com/MrLYC/cmdr"
			var targetPath string

			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .Missing.Field }}")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())
			direct := strategy.NewDirectStrategy()
			direct.SetEnabled(true)
			chain := strategy.NewStrategyChain(direct, rewrite)
			downloadManager.SetStrategyChain(chain)

			fetcher.EXPECT().IsSupport(input).Return(true)
			fetcher.EXPECT().Fetch(name, version, input, gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(os.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())
				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(location).To(Equal(targetPath))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, input)).To(Succeed())
		})

		It("should wrap strategy chain errors", func() {
			chain := strategy.NewStrategyChain(strategy.NewDirectStrategy())
			downloadManager.SetStrategyChain(chain)

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, uri, gomock.Any()).Return(errors.New("timeout")).Times(3)

			_, err := downloadManager.Define(name, version, uri)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to download"))
		})

		DescribeTable("fetch multiple files", func(files map[string]os.FileMode, expected string) {
			var outputDir string

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(name, version, gomock.Any(), gomock.Any()).DoAndReturn(func(name, version, uri, dir string) error {
				outputDir = dir

				for path, mode := range files {
					target := filepath.Join(dir, path)
					Expect(os.MkdirAll(filepath.Dir(target), 0755)).To(Succeed())
					Expect(os.WriteFile(target, []byte(""), mode)).To(Succeed())
				}

				return nil
			})

			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) (core.Command, error) {
				Expect(filepath.Rel(outputDir, location)).To(Equal(expected))
				return nil, nil
			})

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		},
			Entry("single executable even name not match", map[string]os.FileMode{
				"x": 0755,
			}, "x"),
			Entry("perfer to choose by name", map[string]os.FileMode{
				"x":    0755,
				"cmdr": 0644,
			}, "cmdr"),
			Entry("perfer to choose by name when name not match", map[string]os.FileMode{
				"x":  0755,
				"xx": 0644,
			}, "x"),
			Entry("single executable", map[string]os.FileMode{
				"cmdr": 0755,
			}, "cmdr"),
			Entry("single file", map[string]os.FileMode{
				"cmdr": 0644,
			}, "cmdr"),
			Entry("perfer to choose executable", map[string]os.FileMode{
				"cmdr1": 0644,
				"cmdr2": 0755,
			}, "cmdr2"),
			Entry("perfer to choose shorter name", map[string]os.FileMode{
				"cmdr-with-long-name": 0755,
				"cmdr-shortter":       0755,
			}, "cmdr-shortter"),
			Entry("perfer to choose shorter name even if it is not executable", map[string]os.FileMode{
				"cmdr-with-long-name": 0755,
				"cmdr-shortter":       0644,
			}, "cmdr-shortter"),
			Entry("perfer to choose executable even if it is in a sub directory", map[string]os.FileMode{
				"cmdr-with-long-name":               0755,
				"this/a/long/dir/for/cmdr-shortter": 0644,
			}, "this/a/long/dir/for/cmdr-shortter"),
			Entry("perfer to choose the shorter path when multiple files have same name", map[string]os.FileMode{
				"cmdr":         0755,
				"sub/dir/cmdr": 0755,
			}, "cmdr"),
		)

		It("should replace url", func() {
			input := "http://github.com/MrLYC/cmdr"
			replaced := "mock://github.com/MrLYC/cmdr"

			downloadManager.SetReplacements(utils.Replacements{
				{
					Match:    "http://(.*)",
					Template: "mock://{{ index .group 1 }}",
				},
			})

			fetcher.EXPECT().IsSupport(replaced).Return(false)
			baseManager.EXPECT().Define(name, version, replaced)

			Expect(downloadManager.Define(name, version, input)).To(Succeed())
		})
	})

	Context("Factory", func() {
		var (
			cfg     core.Configuration
			rootDir string
		)

		BeforeEach(func() {
			cfg = viper.New()

			var err error
			rootDir, err = os.MkdirTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			cfg.Set(core.CfgKeyCmdrDatabasePath, filepath.Join(rootDir, "cmdr.db"))
		})

		AfterEach(func() {
			os.RemoveAll(rootDir)
		})

		It("should create download manager", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderDownload, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.DownloadManager)
			Expect(ok).To(BeTrue())
		})
	})
})
