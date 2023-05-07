package fetcher_test

import (
	"os"
	"path/filepath"

	. "github.com/mrlyc/cmdr/core/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoGetter", func() {
	var (
		getter *GoGetter
	)

	BeforeEach(func() {
		getter = NewDefaultGoGetter(os.Stdout)
	})

	Context("IsSupport", func() {
		It("should return true", func() {
			Expect(getter.IsSupport(filepath.Join(os.TempDir(), "anything"))).To(BeFalse())
		})

		It("should return false", func() {
			Expect(getter.IsSupport("http://example.com")).To(BeTrue())
		})
	})
})
