package operator_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/operator/mock"
)

var _ = Describe("Operator", func() {
	var (
		ctrl                 *gomock.Controller
		operator1, operator2 *mock.MockOperatorer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		operator1 = mock.NewMockOperatorer(ctrl)
		operator1.EXPECT().String().Return("operator1").AnyTimes()

		operator2 = mock.NewMockOperatorer(ctrl)
		operator2.EXPECT().String().Return("operator2").AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Run", func() {
		Context("ordering", func() {
			var histories []string

			BeforeEach(func() {
				histories = make([]string, 0, 4)

				operator1.EXPECT().Run(gomock.Any()).Do(func(ctx context.Context) (context.Context, error) {
					histories = append(histories, "run-operator1")
					return ctx, nil
				}).AnyTimes()
				operator1.EXPECT().Commit(gomock.Any()).Do(func(ctx context.Context) error {
					histories = append(histories, "finish-operator1")
					return nil
				}).AnyTimes()

				operator2.EXPECT().Run(gomock.Any()).Do(func(ctx context.Context) (context.Context, error) {
					histories = append(histories, "run-operator2")
					return ctx, nil
				}).AnyTimes()
				operator2.EXPECT().Commit(gomock.Any()).Do(func(ctx context.Context) error {
					histories = append(histories, "finish-operator2")
					return nil
				}).AnyTimes()
			})

			It("from new", func() {
				runner := NewOperatorRunner(operator1, operator2)
				Expect(runner.Run(context.Background())).To(Succeed())
				Expect(histories).To(Equal([]string{"run-operator1", "run-operator2", "finish-operator2", "finish-operator1"}))
			})

			It("from add", func() {
				runner := NewOperatorRunner()
				runner.Add(operator1, operator2)
				Expect(runner.Run(context.Background())).To(Succeed())
				Expect(histories).To(Equal([]string{"run-operator1", "run-operator2", "finish-operator2", "finish-operator1"}))
			})
		})

		Context("fail", func() {
			It("check rollback", func() {
				operator1.EXPECT().Run(gomock.Any()).Return(context.Background(), nil)
				operator1.EXPECT().Rollback(gomock.Any())
				operator2.EXPECT().Run(gomock.Any()).Return(context.Background(), fmt.Errorf("error"))
				operator3 := mock.NewMockOperatorer(ctrl)

				runner := NewOperatorRunner(operator1, operator2, operator3)
				Expect(runner.Run(context.Background())).NotTo(Succeed())
			})
		})
	})

	Context("Layout", func() {
		It("mixed", func() {
			runner := NewOperatorRunner(operator1)
			runner.Add(operator2)
			Expect(runner.Layout()).To(Equal([]string{"operator1", "operator2"}))
		})
	})
})
