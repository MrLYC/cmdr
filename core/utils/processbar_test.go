package utils_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/utils"
	"github.com/mrlyc/cmdr/core/utils/mock"
)

var _ = Describe("Processbar", func() {
	var (
		ctrl        *gomock.Controller
		bar         *mock.MockprogressBar
		reader      *mock.MockReadCloser
		progressBar *utils.ProgressBar
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		bar = mock.NewMockprogressBar(ctrl)
		reader = mock.NewMockReadCloser(ctrl)

		progressBar = utils.NewProgressBar(reader, bar)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should add read bytes to progressbar", func() {
		reader.EXPECT().Read(gomock.Any()).Return(666, nil)
		bar.EXPECT().Add(666).Return(nil)

		n, err := progressBar.Read([]byte{})
		Expect(err).To(BeNil())
		Expect(n).To(Equal(666))
	})

	It("should close progressbar", func() {
		reader.EXPECT().Close().Return(nil)
		bar.EXPECT().Finish().Return(nil)

		Expect(progressBar.Close()).To(BeNil())
	})
})
