package core_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Step", func() {
	var (
		ctrl         *gomock.Controller
		step1, step2 *mock.MockSteper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		step1 = mock.NewMockSteper(ctrl)
		step1.EXPECT().String().Return("step1").AnyTimes()

		step2 = mock.NewMockSteper(ctrl)
		step2.EXPECT().String().Return("step2").AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Run", func() {
		Context("ordering", func() {
			var histories []string

			BeforeEach(func() {
				histories = make([]string, 0, 4)

				step1.EXPECT().Run(gomock.Any()).Do(func(ctx context.Context) (context.Context, error) {
					histories = append(histories, "run-step1")
					return ctx, nil
				}).AnyTimes()
				step1.EXPECT().Commit(gomock.Any()).Do(func(ctx context.Context) error {
					histories = append(histories, "finish-step1")
					return nil
				}).AnyTimes()

				step2.EXPECT().Run(gomock.Any()).Do(func(ctx context.Context) (context.Context, error) {
					histories = append(histories, "run-step2")
					return ctx, nil
				}).AnyTimes()
				step2.EXPECT().Commit(gomock.Any()).Do(func(ctx context.Context) error {
					histories = append(histories, "finish-step2")
					return nil
				}).AnyTimes()
			})

			It("from new", func() {
				runner := NewStepRunner(step1, step2)
				Expect(runner.Run(context.Background())).To(Succeed())
				Expect(histories).To(Equal([]string{"run-step1", "run-step2", "finish-step2", "finish-step1"}))
			})

			It("from add", func() {
				runner := NewStepRunner()
				runner.Add(step1, step2)
				Expect(runner.Run(context.Background())).To(Succeed())
				Expect(histories).To(Equal([]string{"run-step1", "run-step2", "finish-step2", "finish-step1"}))
			})
		})

		Context("fail", func() {
			It("check rollback", func() {
				step1.EXPECT().Run(gomock.Any()).Return(context.Background(), nil)
				step1.EXPECT().Rollback(gomock.Any())
				step2.EXPECT().Run(gomock.Any()).Return(context.Background(), fmt.Errorf("error"))
				step3 := mock.NewMockSteper(ctrl)

				runner := NewStepRunner(step1, step2, step3)
				Expect(runner.Run(context.Background())).NotTo(Succeed())
			})
		})
	})

	Context("Layout", func() {
		It("mixed", func() {
			runner := NewStepRunner(step1)
			runner.Add(step2)
			Expect(runner.Layout()).To(Equal([]string{"step1", "step2"}))
		})
	})
})
