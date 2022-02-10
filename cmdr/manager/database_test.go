package manager_test

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmdr"
	"github.com/mrlyc/cmdr/cmdr/manager"
	"github.com/mrlyc/cmdr/cmdr/manager/mock"
)

var _ = Describe("Database", func() {
	var (
		ctrl    *gomock.Controller
		db      *mock.MockDBClient
		dbQuery *mock.MockQuery
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDBClient(ctrl)
		dbQuery = mock.NewMockQuery(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
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
				q.Eq("NameField", "name"),
				q.Eq("VersionField", "version"),
				q.Eq("ActivatedField", true),
				q.Eq("LocationField", "location"),
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

	Context("DatabaseManager", func() {
		var (
			mgr           *manager.DatabaseManager
			commandName   = "command"
			version       = "1.0.0"
			location      = "location"
			existsCommand = manager.Command{
				IDField: 1,
			}
		)

		makeCommandNotFound := func() {
			db.EXPECT().
				Select(
					q.Eq("NameField", commandName),
					q.Eq("VersionField", version),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				First(gomock.Any()).
				Return(storm.ErrNotFound)
		}

		makeActivatedCommandNotFound := func() {
			db.EXPECT().
				Select(
					q.Eq("NameField", commandName),
					q.Eq("ActivatedField", true),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				Find(gomock.Any()).
				Return(storm.ErrNotFound)
		}

		makeCommandFound := func() {
			db.EXPECT().
				Select(
					q.Eq("NameField", commandName),
					q.Eq("VersionField", version),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				First(gomock.Any()).
				DoAndReturn(func(target interface{}) error {
					*target.(*manager.Command) = existsCommand
					return nil
				})
		}

		makeActivatedCommandFound := func() {
			db.EXPECT().
				Select(
					q.Eq("NameField", commandName),
					q.Eq("ActivatedField", true),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				Find(gomock.Any()).
				DoAndReturn(func(target interface{}) error {
					*target.(*[]*manager.Command) = []*manager.Command{&existsCommand}
					return nil
				})
		}

		BeforeEach(func() {
			mgr = manager.NewDatabaseManager(db)
		})

		It("should init database", func() {
			db.EXPECT().Init(gomock.Any()).DoAndReturn(func(data interface{}) error {
				_, ok := data.(*manager.Command)
				Expect(ok).To(BeTrue())
				return nil
			})

			db.EXPECT().ReIndex(gomock.Any()).DoAndReturn(func(data interface{}) error {
				_, ok := data.(*manager.Command)
				Expect(ok).To(BeTrue())
				return nil
			})

			Expect(mgr.Init()).To(Succeed())
		})

		It("should close database", func() {
			db.EXPECT().Close()

			Expect(mgr.Close()).To(Succeed())
		})

		It("should return provider", func() {
			Expect(mgr.Provider()).To(Equal(cmdr.CommandProviderDatabase))
		})

		Context("Define", func() {
			It("should create a command", func() {
				makeCommandNotFound()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.NameField).To(Equal(commandName))
					Expect(command.VersionField).To(Equal(version))
					Expect(command.LocationField).To(Equal(location))
					return nil
				})

				Expect(mgr.Define(commandName, version, location)).To(Succeed())
			})

			It("should update a command", func() {
				makeCommandFound()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.IDField).To(Equal(existsCommand.IDField))
					Expect(command.NameField).To(Equal(commandName))
					Expect(command.VersionField).To(Equal(version))
					Expect(command.LocationField).To(Equal(location))
					return nil
				})

				Expect(mgr.Define(commandName, version, location)).To(Succeed())
			})
		})

		Context("Undefine", func() {
			It("should undefine a command", func() {
				makeCommandFound()
				db.EXPECT().DeleteStruct(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.NameField).To(Equal(commandName))
					Expect(command.VersionField).To(Equal(version))
					return nil
				})

				Expect(mgr.Undefine(commandName, version)).To(Succeed())
			})

			It("should undefine a non-exists command", func() {
				makeCommandNotFound()

				Expect(mgr.Undefine(commandName, version)).To(Succeed())
			})

			It("should not undefine a activated command", func() {
				existsCommand.ActivatedField = true
				makeCommandFound()

				Expect(mgr.Undefine(commandName, version)).NotTo(Succeed())
			})
		})

		Context("Activate", func() {
			It("should activate a command", func() {
				makeCommandFound()
				makeActivatedCommandNotFound()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.IDField).To(Equal(existsCommand.IDField))
					Expect(command.NameField).To(Equal(commandName))
					Expect(command.VersionField).To(Equal(version))
					Expect(command.ActivatedField).To(BeTrue())

					return nil
				})

				Expect(mgr.Activate(commandName, version)).To(Succeed())
			})

			It("should reactivate a command", func() {
				makeCommandFound()
				makeActivatedCommandFound()

				db.EXPECT().Update(gomock.Any()).Return(nil).AnyTimes()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.IDField).To(Equal(existsCommand.IDField))
					Expect(command.NameField).To(Equal(commandName))
					Expect(command.VersionField).To(Equal(version))
					Expect(command.ActivatedField).To(BeTrue())

					return nil
				})

				Expect(mgr.Activate(commandName, version)).To(Succeed())
			})

			It("should return an error because command not found", func() {
				makeCommandNotFound()

				Expect(mgr.Activate(commandName, version)).NotTo(Succeed())
			})
		})

		Context("Deactivate", func() {
			It("should deactivate a command", func() {
				makeActivatedCommandFound()

				db.EXPECT().Update(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command).To(Equal(&existsCommand))
					Expect(command.ActivatedField).To(BeFalse())

					return nil
				})

				Expect(mgr.Deactivate(commandName)).To(Succeed())
			})

			It("should deactivate a non-exists command", func() {
				makeActivatedCommandNotFound()

				Expect(mgr.Deactivate(commandName)).To(Succeed())
			})

			It("should deactivate multiple commands", func() {
				var command1, command2 manager.Command

				db.EXPECT().
					Select(q.Eq("NameField", commandName), q.Eq("ActivatedField", true)).
					Return(dbQuery)

				dbQuery.EXPECT().
					Find(gomock.Any()).
					DoAndReturn(func(target interface{}) error {
						*target.(*[]*manager.Command) = []*manager.Command{
							&command1, &command2,
						}
						return nil
					})

				db.EXPECT().Update(&command1).Return(nil)
				db.EXPECT().Update(&command2).Return(nil)

				Expect(mgr.Deactivate(commandName)).To(Succeed())
			})
		})
	})
})
