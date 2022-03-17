package utils_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Download", func() {
	var (
		downloader *utils.Downloader
	)

	BeforeEach(func() {
		downloader = utils.NewDefaultDownloader(os.Stdout)
	})

	Context("IsSupport", func() {
		It("should return true", func() {
			Expect(downloader.IsSupport(filepath.Join(os.TempDir(), "anything"))).To(BeFalse())
		})

		It("should return false", func() {
			Expect(downloader.IsSupport("http://example.com")).To(BeTrue())
		})
	})
})
