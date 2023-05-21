package fetcher_test

import (
	. "github.com/mrlyc/cmdr/core/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoGetter", func() {
	var (
		installer *GoInstaller
	)

	BeforeEach(func() {
		installer = NewDefaultGoInstaller()
	})

	Context("IsSupport", func() {
		It("should return true", func() {
			Expect(installer.IsSupport("go://github.com/mrlyc/cmdr")).To(BeTrue())
		})
	})

})
