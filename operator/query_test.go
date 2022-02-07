package operator_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/model"
	"github.com/mrlyc/cmdr/define"
	. "github.com/mrlyc/cmdr/operator"
)

var _ = Describe("Query", func() {
	var (
		dbPath string
		db     define.DBClient
	)

	BeforeEach(func() {
		tempDir, err := afero.TempDir(afero.NewOsFs(), "", "")
		Expect(err).To(BeNil())

		dbPath = filepath.Join(tempDir, "cmdr.db")
		db, err = core.NewDBClient(dbPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		db.Close()
		Expect(os.RemoveAll(dbPath)).To(BeNil())
	})

	Context("Run", func() {
		var (
			ctx                context.Context
			commandA, commandB *model.Command
		)

		BeforeEach(func() {
			ctx = context.WithValue(context.Background(), define.ContextKeyDBClient, db)
			commandA = &model.Command{
				Name:      "a",
				Activated: true,
			}
			commandB = &model.Command{
				Name:      "b",
				Activated: true,
			}

			Expect(db.Save(commandA)).To(BeNil())
			Expect(db.Save(commandB)).To(BeNil())
		})

		It("query single command", func() {
			querier := NewCommandsQuerier([]q.Matcher{
				q.Eq("Name", "a"),
			})

			result, err := querier.Run(ctx)
			Expect(err).To(BeNil())
			Expect(result.Value(define.ContextKeyCommands)).To(Equal([]*model.Command{commandA}))
		})

		It("query commands", func() {
			querier := NewCommandsQuerier([]q.Matcher{
				q.Eq("Activated", true),
			})

			result, err := querier.Run(ctx)
			Expect(err).To(BeNil())
			Expect(result.Value(define.ContextKeyCommands)).To(Equal([]*model.Command{commandA, commandB}))
		})

		It("query not found", func() {
			querier := NewCommandsQuerier([]q.Matcher{
				q.Eq("Name", "c"),
			})

			result, err := querier.Run(ctx)
			Expect(err).To(BeNil())
			Expect(result.Value(define.ContextKeyCommands)).To(BeEmpty())
		})
	})
})
