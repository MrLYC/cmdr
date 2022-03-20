package manager_test

import (
	"github.com/asdine/storm/v3/q"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Database", func() {
	var (
		ctrl    *gomock.Controller
		db      *mock.MockDatabase
		dbQuery *mock.MockQuery
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDatabase(ctrl)
		dbQuery = mock.NewMockQuery(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CommandFilter", func() {
		var (
			filter             *manager.CommandFilter
			commandA, commandB *manager.Command
		)

		BeforeEach(func() {
			commandA = &manager.Command{
				Name:      "command-a",
				Version:   "1.0.0",
				Activated: true,
				Location:  "location-a",
			}
			commandB = &manager.Command{
				Name:      "command-b",
				Version:   "1.0.1",
				Activated: false,
				Location:  "location-b",
			}
			filter = manager.NewCommandFilter([]*manager.Command{commandA, commandB})
		})

		It("should return array", func() {
			result, err := filter.All()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]core.Command{commandA, commandB}))
		})

		It("should return error", func() {
			command, err := filter.WithName(commandA.Name).One()
			Expect(err).To(BeNil())
			Expect(command).To(Equal(commandA))
		})

		It("should return 0", func() {
			result, err := filter.Count()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(2))
		})
	})

	Context("CommandQuery", func() {
		var query *manager.CommandQuery

		BeforeEach(func() {
			query = manager.NewCommandQuery(db)
		})

		It("should append matchers", func() {
			query.
				WithName("name").
				WithVersion("version").
				WithActivated(true).
				WithLocation("location")

			db.EXPECT().Select(
				q.Eq("Name", "name"),
				q.Eq("Version", "version"),
				q.Eq("Activated", true),
				q.Eq("Location", "location"),
			).Return(dbQuery)

			Expect(query.Done()).To(Equal(dbQuery))
		})

		It("should return array", func() {
			commands := []*manager.Command{{}}
			db.EXPECT().Select().Return(dbQuery)
			dbQuery.EXPECT().Find(gomock.Any()).DoAndReturn(func(target interface{}) error {
				*target.(*[]*manager.Command) = commands
				return nil
			})

			results, err := query.All()
			Expect(err).To(BeNil())
			Expect(results[0]).To(Equal(commands[0]))
		})

		It("should return command", func() {
			command := &manager.Command{}
			db.EXPECT().Select().Return(dbQuery)
			dbQuery.EXPECT().First(gomock.Any()).DoAndReturn(func(target interface{}) error {
				*target.(*manager.Command) = *command
				return nil
			})

			result, err := query.One()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(command))
		})

		It("should return count", func() {
			db.EXPECT().Select().Return(dbQuery)
			dbQuery.EXPECT().Count(gomock.Any()).Return(1, nil)

			count, err := query.Count()
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})
	})
})
