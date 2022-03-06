package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/manager/mock"
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

	Context("DatabaseManager", func() {
		var (
			mgr           *manager.DatabaseManager
			commandName   = "command"
			version       = "1.0.0"
			location      = "location"
			existsCommand = manager.Command{
				ID: 1,
			}
		)

		makeCommandNotFound := func() {
			db.EXPECT().
				Select(
					q.Eq("Name", commandName),
					q.Eq("Version", version),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				First(gomock.Any()).
				Return(storm.ErrNotFound)
		}

		makeActivatedCommandNotFound := func() {
			db.EXPECT().
				Select(
					q.Eq("Name", commandName),
					q.Eq("Activated", true),
				).
				Return(dbQuery)

			dbQuery.EXPECT().
				Find(gomock.Any()).
				Return(storm.ErrNotFound)
		}

		makeCommandFound := func() {
			db.EXPECT().
				Select(
					q.Eq("Name", commandName),
					q.Eq("Version", version),
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
					q.Eq("Name", commandName),
					q.Eq("Activated", true),
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

		It("should close database", func() {
			db.EXPECT().Close()

			Expect(mgr.Close()).To(Succeed())
		})

		It("should return provider", func() {
			Expect(mgr.Provider()).To(Equal(core.CommandProviderDatabase))
		})

		Context("Define", func() {
			It("should create a command", func() {
				makeCommandNotFound()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Location).To(Equal(location))
					return nil
				})

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(BeNil())
			})

			It("should update a command", func() {
				makeCommandFound()

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.ID).To(Equal(existsCommand.ID))
					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Location).To(Equal(location))
					return nil
				})

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(BeNil())
			})
		})

		Context("Undefine", func() {
			It("should undefine a command", func() {
				makeCommandFound()
				db.EXPECT().DeleteStruct(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					return nil
				})

				Expect(mgr.Undefine(commandName, version)).To(Succeed())
			})

			It("should undefine a non-exists command", func() {
				makeCommandNotFound()

				Expect(mgr.Undefine(commandName, version)).To(Succeed())
			})

			It("should not undefine a activated command", func() {
				existsCommand.Activated = true
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

					Expect(command.ID).To(Equal(existsCommand.ID))
					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Activated).To(BeTrue())

					return nil
				})

				Expect(mgr.Activate(commandName, version)).To(Succeed())
			})

			It("should reactivate a command", func() {
				makeCommandFound()
				makeActivatedCommandFound()

				db.EXPECT().Save(&existsCommand).Return(nil)

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.ID).To(Equal(existsCommand.ID))
					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Activated).To(BeTrue())

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

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command).To(Equal(&existsCommand))
					Expect(command.Activated).To(BeFalse())

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
					Select(q.Eq("Name", commandName), q.Eq("Activated", true)).
					Return(dbQuery)

				dbQuery.EXPECT().
					Find(gomock.Any()).
					DoAndReturn(func(target interface{}) error {
						*target.(*[]*manager.Command) = []*manager.Command{
							&command1, &command2,
						}
						return nil
					})

				db.EXPECT().Save(&command1).Return(nil)
				db.EXPECT().Save(&command2).Return(nil)

				Expect(mgr.Deactivate(commandName)).To(Succeed())
			})
		})
	})

	Context("Factory", func() {
		var (
			cfg     core.Configuration
			rootDir string
		)

		BeforeEach(func() {
			cfg = viper.New()

			var err error
			rootDir, err = os.MkdirTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			cfg.Set(core.CfgKeyCmdrDatabasePath, filepath.Join(rootDir, "cmdr.db"))
		})

		AfterEach(func() {
			os.RemoveAll(rootDir)
		})

		It("should new a manager", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.DatabaseManager)
			Expect(ok).To(BeTrue())
		})

		It("should new a initializer", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.DatabaseManager)
			Expect(ok).To(BeTrue())
		})
	})

	Context("DatabaseMigrator", func() {
		var (
			migrator *manager.DatabaseMigrator
		)

		BeforeEach(func() {
			migrator = manager.NewDatabaseMigrator(func() (manager.DBClient, error) {
				return db, nil
			})
		})

		It("should migrate models", func() {
			histories := make(map[string][]string)

			db.EXPECT().Init(gomock.Any()).DoAndReturn(func(data interface{}) error {
				key := fmt.Sprintf("%T", data)
				histories[key] = append(histories[key], "Init")
				return nil
			}).AnyTimes()
			db.EXPECT().ReIndex(gomock.Any()).DoAndReturn(func(data interface{}) error {
				key := fmt.Sprintf("%T", data)
				histories[key] = append(histories[key], "ReIndex")
				return nil
			}).AnyTimes()
			db.EXPECT().Close().Return(nil)

			Expect(migrator.Init()).To(Succeed())
			Expect(histories["*manager.Command"]).To(Equal([]string{"Init", "ReIndex"}))
		})
	})
})
