package operator_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
)

var _ = Describe("Checker", func() {
	Context("BinariesChecker", func() {
		var (
			ctx      context.Context
			location string
			command  *model.Command
		)

		BeforeEach(func() {
			binary, err := afero.TempFile(define.FS, "", "")
			Expect(err).To(BeNil())

			location = binary.Name()
			command = &model.Command{
				Location: location,
			}
			ctx = context.WithValue(
				context.Background(),
				define.ContextKeyCommands,
				[]*model.Command{command},
			)
		})

		AfterEach(func() {
			Expect(define.FS.Remove(location)).To(Succeed())
		})

		It("should success", func() {
			checker := operator.NewBinariesChecker()

			_, err := checker.Run(ctx)
			Expect(err).To(BeNil())
		})

		It("should fail because binary not exists", func() {
			command.Location = "not_exists"
			checker := operator.NewBinariesChecker()

			_, err := checker.Run(ctx)
			Expect(err).NotTo(BeNil())
		})

		It("should fail because context not found", func() {
			checker := operator.NewBinariesChecker()

			_, err := checker.Run(context.Background())
			Expect(err).NotTo(BeNil())
		})
	})

	Context("CommandsChecker", func() {
		It("should fail because context not found", func() {
			checker := operator.NewCommandsChecker()

			_, err := checker.Run(context.Background())
			Expect(err).NotTo(BeNil())
		})

		It("should fail because commands is nil", func() {
			checker := operator.NewCommandsChecker()

			_, err := checker.Run(context.WithValue(
				context.Background(),
				define.ContextKeyCommands,
				nil,
			))
			Expect(err).NotTo(BeNil())
		})

		It("should fail because commands is empty", func() {
			checker := operator.NewCommandsChecker()

			_, err := checker.Run(context.WithValue(
				context.Background(),
				define.ContextKeyCommands,
				[]*model.Command{},
			))
			Expect(err).NotTo(BeNil())
		})

		It("should success", func() {
			checker := operator.NewCommandsChecker()

			_, err := checker.Run(context.WithValue(
				context.Background(),
				define.ContextKeyCommands,
				[]*model.Command{{}},
			))
			Expect(err).To(BeNil())
		})
	})
})
