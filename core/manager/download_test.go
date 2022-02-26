package manager_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Download", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
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
			downloadManager = manager.NewDownloadManager(baseManager, fetcher)
		})

		It("should call base manager", func() {
			fetcher.EXPECT().IsSupport(uri).Return(false)
			baseManager.EXPECT().Define(name, version, uri)

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		})

		It("should call with downloaded file", func() {
			var targetPath string

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any()).DoAndReturn(func(uri, dir string) error {
				targetPath = filepath.Join(dir, "cmdr")
				Expect(ioutil.WriteFile(targetPath, []byte(""), 0755)).To(Succeed())

				return nil
			})
			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) error {
				Expect(targetPath).To(Equal(location))
				return nil
			})

			Expect(downloadManager.Define(name, version, uri)).To(Succeed())
		})

		DescribeTable("fetch multiple files", func(files map[string]os.FileMode, expected string) {
			var outputDir string

			fetcher.EXPECT().IsSupport(uri).Return(true)
			fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any()).DoAndReturn(func(uri, dir string) error {
				outputDir = dir

				for path, mode := range files {
					target := filepath.Join(dir, path)
					Expect(os.MkdirAll(filepath.Dir(target), 0755)).To(Succeed())
					Expect(ioutil.WriteFile(target, []byte(""), mode)).To(Succeed())
				}

				return nil
			})

			baseManager.EXPECT().Define(name, version, gomock.Any()).DoAndReturn(func(name, version, location string) error {
				Expect(filepath.Rel(outputDir, location)).To(Equal(expected))
				return nil
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
	})
})
