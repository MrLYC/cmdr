package manager

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/strategy"
)

type internalFetcher struct {
	seenURI string
	err     error
}

func (f *internalFetcher) IsSupport(string) bool { return true }
func (f *internalFetcher) Fetch(name, version, uri, dst string) error {
	f.seenURI = uri
	if f.err != nil {
		return f.err
	}
	return os.WriteFile(filepath.Join(dst, name), []byte("binary"), 0755)
}

func managerTempDir() string {
	dir, err := os.MkdirTemp("", "cmdr-manager-test-*")
	Expect(err).NotTo(HaveOccurred())
	return dir
}

var _ = Describe("Manager internals", func() {
	DescribeTable("resolveVersionFromLocation",
		func(name, fallback, location, expected string) {
			Expect(resolveVersionFromLocation(name, fallback, location)).To(Equal(expected))
		},
		Entry("extracts the matching command suffix", "cmdr", "1", "/tmp/cmdr_1.2.3", "1.2.3"),
		Entry("keeps fallback for another command", "cmdr", "1", "/tmp/other_1.2.3", "1"),
		Entry("keeps fallback for an empty suffix", "cmdr", "1", "/tmp/cmdr_", "1"),
	)

	Context("download manager", func() {
		It("returns search errors for missing paths", func() {
			dir := managerTempDir()
			defer os.RemoveAll(dir)

			manager := NewDownloadManager(nil, nil, 1, nil)
			_, err := manager.search("cmdr", filepath.Join(dir, "missing"))
			Expect(err).To(HaveOccurred())
		})

		It("rewrites URIs before fetching", func() {
			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .URI }}?mirror=1")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())

			fetcher := &internalFetcher{}
			manager := NewDownloadManager(nil, nil, 1, nil)
			manager.SetStrategyChain(strategy.NewStrategyChain(rewrite))
			dst := managerTempDir()
			defer os.RemoveAll(dst)

			result, err := manager.fetch(fetcher, "cmdr", "1.0.0", "https://example.com/cmdr", dst)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetcher.seenURI).To(Equal("https://example.com/cmdr?mirror=1"))
			Expect(result).To(Equal(filepath.Join(dst, "cmdr")))
		})

		It("returns fetcher errors", func() {
			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ .URI }}")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())

			fetcher := &internalFetcher{err: errors.New("fetch failed")}
			manager := NewDownloadManager(nil, nil, 1, nil)
			manager.SetStrategyChain(strategy.NewStrategyChain(rewrite))
			dst := managerTempDir()
			defer os.RemoveAll(dst)

			_, err := manager.fetch(fetcher, "cmdr", "1.0.0", "https://example.com/cmdr", dst)
			Expect(err).To(HaveOccurred())
		})

		It("keeps prepared URIs when direct strategy wins", func() {
			cfg := viper.New()
			cfg.Set(core.CfgKeyDownloadRewriteRule, "{{ call .URI }}")
			rewrite := strategy.NewRewriteStrategy()
			Expect(rewrite.Configure(cfg)).To(Succeed())

			direct := strategy.NewDirectStrategy()
			direct.SetEnabled(true)
			fetcher := &internalFetcher{}
			manager := NewDownloadManager(nil, nil, 1, nil)
			manager.SetStrategyChain(strategy.NewStrategyChain(direct, rewrite))
			manager.SetReplacements(nil)
			dst := managerTempDir()
			defer os.RemoveAll(dst)

			_, err := manager.fetch(fetcher, "cmdr", "1.0.0", "https://example.com/cmdr", dst)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetcher.seenURI).To(Equal("https://example.com/cmdr"))
		})
	})

	Context("binary manager", func() {
		It("creates configured managers", func() {
			cfg := viper.New()
			cfg.Set(core.CfgKeyCmdrBinDir, "bin")
			cfg.Set(core.CfgKeyCmdrShimsDir, "shims")
			cfg.Set(core.CfgKeyCmdrLinkMode, "link")
			Expect(newBinaryManagerByConfiguration(cfg)).NotTo(BeNil())

			cfg.Set(core.CfgKeyCmdrLinkMode, "copy")
			Expect(NewBinaryManagerWithCopy("bin", "shims", 0755)).NotTo(BeNil())
			Expect(newBinaryManagerByConfiguration(cfg)).NotTo(BeNil())
		})

		It("returns query errors when shim paths are missing", func() {
			binDir := managerTempDir()
			defer os.RemoveAll(binDir)
			shimRoot := managerTempDir()
			defer os.RemoveAll(shimRoot)

			mgr := NewBinaryManagerWithLink(filepath.Join(binDir, "bin"), filepath.Join(shimRoot, "missing"), 0755)
			_, err := mgr.Query()
			Expect(err).To(HaveOccurred())
		})
	})
})
