package initializer_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/initializer"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Chaining", func() {
	var (
		ctrl         *gomock.Controller
		initializer1 *mock.MockInitializer
		initializer2 *mock.MockInitializer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		initializer1 = mock.NewMockInitializer(ctrl)
		initializer2 = mock.NewMockInitializer(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should new chaining with initializers", func() {
		chaining := initializer.NewChaining(initializer1, initializer2)
		Expect(chaining.GetInitializers()).To(HaveLen(2))
	})

	It("should add initializers", func() {
		chaining := initializer.NewChaining()
		Expect(chaining.GetInitializers()).To(BeEmpty())

		chaining.Add(initializer1, nil, 1, initializer2)
		Expect(chaining.GetInitializers()).To(HaveLen(2))
	})

	It("should init all initializers", func() {
		initializer1.EXPECT().Init()
		initializer2.EXPECT().Init()

		chaining := initializer.NewChaining(initializer1, initializer2)
		Expect(chaining.Init()).To(Succeed())
	})

	It("should init all initializers even some of them fail", func() {
		initializer1.EXPECT().Init().Return(fmt.Errorf("testing"))
		initializer2.EXPECT().Init()

		chaining := initializer.NewChaining(initializer1, initializer2)
		Expect(chaining.Init()).NotTo(Succeed())
	})
})
