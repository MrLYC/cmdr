package fetcher_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	g "github.com/hashicorp/go-getter"
	. "github.com/mrlyc/cmdr/core/fetcher"
	"github.com/mrlyc/cmdr/core/fetcher/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoGetter", func() {
	var (
		ctrl     *gomock.Controller
		getter   *GoGetter
		detector *mock.MockDetector
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		detector = mock.NewMockDetector(ctrl)

		getter = NewGoGetter(nil, []g.Detector{detector}, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("IsSupport", func() {
		It("should return false", func() {
			detector.EXPECT().Detect(gomock.Any(), gomock.Any()).Return("", false, nil)

			Expect(getter.IsSupport("testing")).To(BeFalse())
		})

		It("should return true", func() {
			detector.EXPECT().Detect(gomock.Any(), gomock.Any()).Return("", true, nil)

			Expect(getter.IsSupport("testing")).To(BeTrue())
		})
	})

	Context("Fetch", func() {
		It("should return a error", func() {
			detector.EXPECT().Detect(gomock.Any(), gomock.Any()).Return("", false, fmt.Errorf("testing"))

			Expect(getter.Fetch("name", "version", "uri", "dst")).To(HaveOccurred())
		})
	})
})
