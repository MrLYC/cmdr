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
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Database", func() {
	var (
		ctrl      *gomock.Controller
		db        *mock.MockDatabase
		dbQuery   *mock.MockQuery
		binaryMgr *mock.MockCommandManager
		dbFactory func() (core.Database, error)
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDatabase(ctrl)
		dbQuery = mock.NewMockQuery(ctrl)
		binaryMgr = mock.NewMockCommandManager(ctrl)
		core.SetDatabaseFactory(func() (core.Database, error) {
			return db, nil
		})
		dbFactory = core.GetDatabaseFactory()
	})

	AfterEach(func() {
		ctrl.Finish()
		core.SetDatabaseFactory(dbFactory)
	})

	Context("DatabaseManager", func() {
		var (
			mgr           *manager.DatabaseManager
			commandName   = "command"
			version       = "1.0"
			location      = "location"
			existsCommand = manager.Command{
				ID: 1,
			}
		)

		makeCommandNotFound := func() {
			db.EXPECT().
				Select(
					q.Eq("Name", commandName),
					q.Or(
						q.Eq("Version", "1.0"),
						q.Eq("Version", "1.0.0"),
					),
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
					q.Or(
						q.Eq("Version", "1.0"),
						q.Eq("Version", "1.0.0"),
					),
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
			existsCommand = manager.Command{ID: 1}
			mgr = manager.NewDatabaseManager(db, binaryMgr)
		})

		It("should close database", func() {
			binaryMgr.EXPECT().Close()
			db.EXPECT().Close()

			Expect(mgr.Close()).To(Succeed())
		})

		It("should return close manager errors", func() {
			binaryMgr.EXPECT().Close().Return(fmt.Errorf("close failed"))
			binaryMgr.EXPECT().Provider().Return(core.CommandProviderBinary)

			Expect(mgr.Close()).To(HaveOccurred())
		})

		It("should return provider", func() {
			Expect(mgr.Provider()).To(Equal(core.CommandProviderDatabase))
		})

		It("should return query", func() {
			query, err := mgr.Query()
			Expect(err).NotTo(HaveOccurred())
			Expect(query).NotTo(BeNil())
		})

		Context("Define", func() {
			var (
				binaryCommand *mock.MockCommand
			)

			BeforeEach(func() {
				binaryCommand = mock.NewMockCommand(ctrl)
				binaryCommand.EXPECT().GetLocation().Return("binary").AnyTimes()
			})

			It("should create a command", func() {
				makeCommandNotFound()

				binaryMgr.EXPECT().Define(commandName, version, location).Return(binaryCommand, nil)

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Location).To(Equal(binaryCommand.GetLocation()))
					return nil
				})

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(BeNil())
			})

			It("should update a command", func() {
				makeCommandFound()

				binaryMgr.EXPECT().Define(commandName, version, location).Return(binaryCommand, nil)

				db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
					command, ok := data.(*manager.Command)
					Expect(ok).To(BeTrue())

					Expect(command.ID).To(Equal(existsCommand.ID))
					Expect(command.Name).To(Equal(commandName))
					Expect(command.Version).To(Equal(version))
					Expect(command.Location).To(Equal(binaryCommand.GetLocation()))
					return nil
				})

				result, err := mgr.Define(commandName, version, location)
				Expect(err).To(BeNil())
				Expect(result.GetLocation()).To(Equal(binaryCommand.GetLocation()))
			})

			It("should return binary manager errors", func() {
				binaryMgr.EXPECT().Define(commandName, version, location).Return(nil, fmt.Errorf("define failed"))

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(MatchError("define failed"))
			})

			It("should return query errors", func() {
				binaryMgr.EXPECT().Define(commandName, version, location).Return(binaryCommand, nil)
				db.EXPECT().
					Select(
						q.Eq("Name", commandName),
						q.Or(
							q.Eq("Version", "1.0"),
							q.Eq("Version", "1.0.0"),
						),
					).
					Return(dbQuery)
				dbQuery.EXPECT().First(gomock.Any()).Return(fmt.Errorf("select failed"))

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("define command failed"))
			})

			It("should return save errors", func() {
				makeCommandNotFound()
				binaryMgr.EXPECT().Define(commandName, version, location).Return(binaryCommand, nil)
				db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("save failed"))

				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("save command failed"))
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

				binaryMgr.EXPECT().Undefine(commandName, version).Return(nil)

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

			It("should return query errors", func() {
				db.EXPECT().
					Select(
						q.Eq("Name", commandName),
						q.Or(
							q.Eq("Version", "1.0"),
							q.Eq("Version", "1.0.0"),
						),
					).
					Return(dbQuery)
				dbQuery.EXPECT().First(gomock.Any()).Return(fmt.Errorf("select failed"))

				Expect(mgr.Undefine(commandName, version)).NotTo(Succeed())
			})

			It("should return delete errors", func() {
				existsCommand.Activated = false
				makeCommandFound()
				db.EXPECT().DeleteStruct(gomock.Any()).Return(fmt.Errorf("delete failed"))

				Expect(mgr.Undefine(commandName, version)).NotTo(Succeed())
			})

			It("should return binary manager errors", func() {
				existsCommand.Activated = false
				makeCommandFound()
				db.EXPECT().DeleteStruct(gomock.Any()).Return(nil)
				binaryMgr.EXPECT().Undefine(commandName, version).Return(fmt.Errorf("undefine failed"))

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

				binaryMgr.EXPECT().Activate(commandName, version).Return(nil)

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

				binaryMgr.EXPECT().Deactivate(commandName).Return(nil)
				binaryMgr.EXPECT().Activate(commandName, version).Return(nil)

				Expect(mgr.Activate(commandName, version)).To(Succeed())
			})

			It("should return an error because command not found", func() {
				makeCommandNotFound()

				Expect(mgr.Activate(commandName, version)).NotTo(Succeed())
			})

			It("should return deactivate errors", func() {
				makeCommandFound()
				makeActivatedCommandFound()
				db.EXPECT().Save(&existsCommand).Return(fmt.Errorf("save failed"))

				Expect(mgr.Activate(commandName, version)).NotTo(Succeed())
			})

			It("should return save errors", func() {
				makeCommandFound()
				makeActivatedCommandNotFound()
				db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("save failed"))

				Expect(mgr.Activate(commandName, version)).NotTo(Succeed())
			})

			It("should return binary manager errors", func() {
				makeCommandFound()
				makeActivatedCommandNotFound()
				db.EXPECT().Save(gomock.Any()).Return(nil)
				binaryMgr.EXPECT().Activate(commandName, version).Return(fmt.Errorf("activate failed"))

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

				binaryMgr.EXPECT().Deactivate(gomock.Any()).Return(nil)

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

				binaryMgr.EXPECT().Deactivate(commandName).Return(nil)

				Expect(mgr.Deactivate(commandName)).To(Succeed())
			})

			It("should return select errors", func() {
				db.EXPECT().
					Select(q.Eq("Name", commandName), q.Eq("Activated", true)).
					Return(dbQuery)
				dbQuery.EXPECT().Find(gomock.Any()).Return(fmt.Errorf("select failed"))

				Expect(mgr.Deactivate(commandName)).NotTo(Succeed())
			})

			It("should return save errors", func() {
				makeActivatedCommandFound()
				db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("save failed"))

				Expect(mgr.Deactivate(commandName)).NotTo(Succeed())
			})

			It("should return binary manager errors", func() {
				makeActivatedCommandFound()
				db.EXPECT().Save(gomock.Any()).Return(nil)
				binaryMgr.EXPECT().Deactivate(commandName).Return(fmt.Errorf("deactivate failed"))

				Expect(mgr.Deactivate(commandName)).NotTo(Succeed())
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
})
