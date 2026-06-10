package manager_test

import (
	"fmt"

	"github.com/asdine/storm/v3/q"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

	Context("Command", func() {
		DescribeTable("should shorttern command version", func(version, expected string) {
			cmd := manager.Command{
				Version: version,
			}
			Expect(cmd.GetVersion()).To(Equal(expected))
		},
			Entry("1", "1", "1"),
			Entry("1.0", "1.0", "1"),
			Entry("1.0.0", "1.0.0", "1"),
			Entry("1.1", "1.1", "1.1"),
			Entry("1.1.0", "1.1.0", "1.1"),
			Entry("1.1.1", "1.1.1", "1.1.1"),
		)

		It("should render command string", func() {
			Expect((&manager.Command{Name: "cmdr", Version: "1.0.0"}).String()).To(Equal("cmdr(1.0.0)"))
			Expect((&manager.Command{Name: "cmdr", Version: "1.0.0", Activated: true}).String()).To(Equal("*cmdr(1.0.0)"))
		})
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

		It("should return an error for empty one", func() {
			_, err := filter.WithName("missing").One()
			Expect(err).To(HaveOccurred())
		})

		It("should return 0", func() {
			result, err := filter.Count()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(2))
		})

		DescribeTable("return command-a", func(fn func(core.CommandQuery)) {
			fn(filter.WithName(commandA.Name))
			result, err := filter.All()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]core.Command{commandA}))
		},
			Entry("by version 1", func(query core.CommandQuery) {
				query.WithVersion("1")
			}),
			Entry("by version 1.0", func(query core.CommandQuery) {
				query.WithVersion("1.0")
			}),
			Entry("by version 1.0.0", func(query core.CommandQuery) {
				query.WithVersion("1.0.0")
			}),
		)

		DescribeTable("return command-b", func(fn func(core.CommandQuery)) {
			fn(filter.WithName(commandB.Name))
			result, err := filter.All()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]core.Command{commandB}))
		},
			Entry("by version 1.0.1", func(query core.CommandQuery) {
				query.WithVersion("1.0.1")
			}),
		)

		It("should filter by activation and location and add commands", func() {
			result, err := filter.WithActivated(true).All()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]core.Command{commandA}))

			filter = manager.NewCommandFilter([]*manager.Command{commandA, commandB})
			result, err = filter.WithLocation("location-b").All()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]core.Command{commandB}))

			filter = manager.NewCommandFilter(nil)
			filter.AddCommand(commandA, commandB)
			result, err = filter.All()
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(2))
			Expect(result[0].GetName()).To(Equal(commandA.Name))
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
				WithVersion("1.0").
				WithActivated(true).
				WithLocation("location")

			db.EXPECT().Select(
				q.Eq("Name", "name"),
				q.Or(
					q.Eq("Version", "1.0"),
					q.Eq("Version", "1.0.0"),
				),
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

		It("should return query errors", func() {
			db.EXPECT().Select().Return(dbQuery).Times(3)
			dbQuery.EXPECT().Find(gomock.Any()).Return(fmt.Errorf("find failed"))
			_, err := query.All()
			Expect(err).To(MatchError("find failed"))

			query = manager.NewCommandQuery(db)
			dbQuery.EXPECT().First(gomock.Any()).Return(fmt.Errorf("first failed"))
			_, err = query.One()
			Expect(err).To(MatchError("first failed"))

			query = manager.NewCommandQuery(db)
			dbQuery.EXPECT().Count(gomock.Any()).Return(0, fmt.Errorf("count failed"))
			_, err = query.Count()
			Expect(err).To(MatchError("count failed"))
		})
	})
})
