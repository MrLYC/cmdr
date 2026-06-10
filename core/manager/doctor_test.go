package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/mock"
)

var _ = Describe("Doctor", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("DoctorManager", func() {
		var (
			mainQuery, recorderQuery     *mock.MockCommandQuery
			mainManager, recorderManager *mock.MockCommandManager
			doctor                       *manager.DoctorManager
		)

		BeforeEach(func() {
			mainQuery = mock.NewMockCommandQuery(ctrl)
			recorderQuery = mock.NewMockCommandQuery(ctrl)
			mainManager = mock.NewMockCommandManager(ctrl)
			recorderManager = mock.NewMockCommandManager(ctrl)
			doctor = manager.NewDoctorManager(mainManager, recorderManager)
		})

		Context("Commands", func() {
			It("should call both managers in order", func() {
				firstCommand := mock.NewMockCommand(ctrl)
				secondCommand := mock.NewMockCommand(ctrl)
				gomock.InOrder(
					mainManager.EXPECT().Define("cmdr", "1.0.0", "/tmp/cmdr").Return(firstCommand, nil),
					recorderManager.EXPECT().Define("cmdr", "1.0.0", "/tmp/cmdr").Return(secondCommand, nil),
					mainManager.EXPECT().Undefine("cmdr", "1.0.0").Return(nil),
					recorderManager.EXPECT().Undefine("cmdr", "1.0.0").Return(nil),
					mainManager.EXPECT().Activate("cmdr", "1.0.0").Return(nil),
					recorderManager.EXPECT().Activate("cmdr", "1.0.0").Return(nil),
					mainManager.EXPECT().Deactivate("cmdr").Return(nil),
					recorderManager.EXPECT().Deactivate("cmdr").Return(nil),
				)

				command, err := doctor.Define("cmdr", "1.0.0", "/tmp/cmdr")
				Expect(err).NotTo(HaveOccurred())
				Expect(command).To(Equal(secondCommand))
				Expect(doctor.Undefine("cmdr", "1.0.0")).To(Succeed())
				Expect(doctor.Activate("cmdr", "1.0.0")).To(Succeed())
				Expect(doctor.Deactivate("cmdr")).To(Succeed())
				Expect(doctor.Provider()).To(Equal(core.CommandProviderDoctor))
			})

			It("should stop on the first command error", func() {
				mainManager.EXPECT().Define("cmdr", "1.0.0", "/tmp/cmdr").Return(nil, fmt.Errorf("define failed"))
				_, err := doctor.Define("cmdr", "1.0.0", "/tmp/cmdr")
				Expect(err).To(MatchError("define failed"))

				mainManager.EXPECT().Undefine("cmdr", "1.0.0").Return(fmt.Errorf("undefine failed"))
				Expect(doctor.Undefine("cmdr", "1.0.0")).To(MatchError("undefine failed"))

				mainManager.EXPECT().Activate("cmdr", "1.0.0").Return(fmt.Errorf("activate failed"))
				Expect(doctor.Activate("cmdr", "1.0.0")).To(MatchError("activate failed"))

				mainManager.EXPECT().Deactivate("cmdr").Return(fmt.Errorf("deactivate failed"))
				Expect(doctor.Deactivate("cmdr")).To(MatchError("deactivate failed"))
			})

			It("should collect close errors from both managers", func() {
				mainManager.EXPECT().Close().Return(fmt.Errorf("main close failed"))
				recorderManager.EXPECT().Close().Return(fmt.Errorf("recorder close failed"))

				err := doctor.Close()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("main close failed"))
				Expect(err.Error()).To(ContainSubstring("recorder close failed"))
			})
		})

		Context("Query", func() {
			It("should return recoder query directly", func() {
				gomock.InOrder(
					mainManager.EXPECT().Query().Return(nil, fmt.Errorf("main error")),
					recorderManager.EXPECT().Query().Return(recorderQuery, nil),
				)
				query, err := doctor.Query()
				Expect(err).To(BeNil())
				Expect(query).To(Equal(recorderQuery))
			})

			It("should return main query directly", func() {
				gomock.InOrder(
					mainManager.EXPECT().Query().Return(mainQuery, nil),
					recorderManager.EXPECT().Query().Return(nil, fmt.Errorf("recorder error")),
				)
				query, err := doctor.Query()
				Expect(err).To(BeNil())
				Expect(query).To(Equal(mainQuery))
			})

			It("should return commands from main query when recorder query failed", func() {
				command := mock.NewMockCommand(ctrl)
				command.EXPECT().GetName().Return("command").AnyTimes()
				command.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
				command.EXPECT().GetLocation().Return("/path/to/command").AnyTimes()
				command.EXPECT().GetActivated().Return(true).AnyTimes()

				gomock.InOrder(
					mainManager.EXPECT().Query().Return(mainQuery, nil),
					recorderManager.EXPECT().Query().Return(recorderQuery, nil),
					mainQuery.EXPECT().All().Return([]core.Command{command}, nil),
					recorderQuery.EXPECT().All().Return(nil, fmt.Errorf("recorder error")),
				)

				query, err := doctor.Query()
				Expect(err).To(BeNil())

				count, err := query.Count()
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))

				result, err := query.One()
				Expect(err).To(BeNil())

				Expect(result.GetName()).To(Equal("command"))
				Expect(result.GetVersion()).To(Equal("1.0.1"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
			})

			It("should return commands from recorder query when main query failed", func() {
				command := mock.NewMockCommand(ctrl)
				command.EXPECT().GetName().Return("command").AnyTimes()
				command.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
				command.EXPECT().GetLocation().Return("/path/to/command").AnyTimes()
				command.EXPECT().GetActivated().Return(true).AnyTimes()

				gomock.InOrder(
					mainManager.EXPECT().Query().Return(mainQuery, nil),
					recorderManager.EXPECT().Query().Return(recorderQuery, nil),
					mainQuery.EXPECT().All().Return(nil, fmt.Errorf("recorder error")),
					recorderQuery.EXPECT().All().Return([]core.Command{command}, nil),
				)

				query, err := doctor.Query()
				Expect(err).To(BeNil())

				count, err := query.Count()
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))

				result, err := query.One()
				Expect(err).To(BeNil())

				Expect(result.GetName()).To(Equal("command"))
				Expect(result.GetVersion()).To(Equal("1.0.1"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
			})

			It("should merge commands", func() {
				mainCommand := mock.NewMockCommand(ctrl)
				mainCommand.EXPECT().GetName().Return("command").AnyTimes()
				mainCommand.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
				mainCommand.EXPECT().GetLocation().Return("/path/to/command").AnyTimes()
				mainCommand.EXPECT().GetActivated().Return(false).AnyTimes()

				recorderCommand := mock.NewMockCommand(ctrl)
				recorderCommand.EXPECT().GetName().Return("command").AnyTimes()
				recorderCommand.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
				recorderCommand.EXPECT().GetLocation().Return("no_important").AnyTimes()
				recorderCommand.EXPECT().GetActivated().Return(true).AnyTimes()

				gomock.InOrder(
					mainManager.EXPECT().Query().Return(mainQuery, nil),
					recorderManager.EXPECT().Query().Return(recorderQuery, nil),
					mainQuery.EXPECT().All().Return([]core.Command{mainCommand}, nil),
					recorderQuery.EXPECT().All().Return([]core.Command{recorderCommand}, nil),
				)

				query, err := doctor.Query()
				Expect(err).To(BeNil())

				count, err := query.Count()
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))

				result, err := query.One()
				Expect(err).To(BeNil())

				Expect(result.GetName()).To(Equal("command"))
				Expect(result.GetVersion()).To(Equal("1.0.1"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
			})

			It("should handle multiple versions", func() {
				mainCommand1 := mock.NewMockCommand(ctrl)
				mainCommand1.EXPECT().GetName().Return("command").AnyTimes()
				mainCommand1.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
				mainCommand1.EXPECT().GetLocation().Return("/path/to/command1").AnyTimes()
				mainCommand1.EXPECT().GetActivated().Return(true).AnyTimes()

				mainCommand2 := mock.NewMockCommand(ctrl)
				mainCommand2.EXPECT().GetName().Return("command").AnyTimes()
				mainCommand2.EXPECT().GetVersion().Return("2.0.0").AnyTimes()
				mainCommand2.EXPECT().GetLocation().Return("/path/to/command2").AnyTimes()
				mainCommand2.EXPECT().GetActivated().Return(false).AnyTimes()

				recorderCommand1 := mock.NewMockCommand(ctrl)
				recorderCommand1.EXPECT().GetName().Return("command").AnyTimes()
				recorderCommand1.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
				recorderCommand1.EXPECT().GetLocation().Return("no_important1").AnyTimes()
				recorderCommand1.EXPECT().GetActivated().Return(true).AnyTimes()

				recorderCommand2 := mock.NewMockCommand(ctrl)
				recorderCommand2.EXPECT().GetName().Return("command").AnyTimes()
				recorderCommand2.EXPECT().GetVersion().Return("2.0.0").AnyTimes()
				recorderCommand2.EXPECT().GetLocation().Return("no_important2").AnyTimes()
				recorderCommand2.EXPECT().GetActivated().Return(false).AnyTimes()

				gomock.InOrder(
					mainManager.EXPECT().Query().Return(mainQuery, nil),
					recorderManager.EXPECT().Query().Return(recorderQuery, nil),
					mainQuery.EXPECT().All().Return([]core.Command{mainCommand1, mainCommand2}, nil),
					recorderQuery.EXPECT().All().Return([]core.Command{recorderCommand1, recorderCommand2}, nil),
				)

				query, err := doctor.Query()
				Expect(err).To(BeNil())

				count, err := query.Count()
				Expect(err).To(BeNil())
				Expect(count).To(Equal(2))
			})
		})
	})

	Context("CommandDoctor", func() {
		var (
			query   *mock.MockCommandQuery
			mgr     *mock.MockCommandManager
			command *mock.MockCommand
			doctor  *manager.CommandDoctor
			rootDir string
		)

		BeforeEach(func() {
			query = mock.NewMockCommandQuery(ctrl)
			mgr = mock.NewMockCommandManager(ctrl)
			command = mock.NewMockCommand(ctrl)
			doctor = manager.NewCommandDoctor(mgr, "")

			rootDir, _ = os.MkdirTemp("", "")

			command.EXPECT().GetName().Return("command").AnyTimes()
			command.EXPECT().GetVersion().Return("1.0.1").AnyTimes()
			command.EXPECT().GetLocation().Return(filepath.Join(rootDir, "command")).AnyTimes()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
		})

		It("should return error when query failed", func() {
			mgr.EXPECT().Query().Return(nil, fmt.Errorf("error"))

			Expect(doctor.Fix(false)).NotTo(BeNil())
		})

		It("should return error when get commands failed", func() {
			mgr.EXPECT().Query().Return(query, nil)
			query.EXPECT().All().Return(nil, fmt.Errorf("error"))

			Expect(doctor.Fix(false)).NotTo(BeNil())
		})

		Context("Command not available", func() {
			BeforeEach(func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{command}, nil)
			})

			It("should remove activated command", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Deactivate(command.GetName())
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion())

				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should continue when deactivate and undefine fail", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Deactivate(command.GetName()).Return(fmt.Errorf("deactivate failed"))
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion()).Return(fmt.Errorf("undefine failed"))

				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should not mutate commands in dry-run mode", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()

				Expect(doctor.Fix(true)).To(Succeed())
			})

			It("should remove non-activate command", func() {
				command.EXPECT().GetActivated().Return(false).AnyTimes()
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion())

				Expect(doctor.Fix(false)).To(Succeed())
			})
		})

		Context("Command available", func() {
			BeforeEach(func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{command}, nil)
				Expect(os.WriteFile(command.GetLocation(), []byte("command"), 0755)).To(Succeed())
			})

			It("should re-activate activated command without re-defining", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Activate(command.GetName(), command.GetVersion())

				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should continue when re-activate fails", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Activate(command.GetName(), command.GetVersion()).Return(fmt.Errorf("activate failed"))

				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should log re-activation in dry-run mode", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()

				Expect(doctor.Fix(true)).To(Succeed())
			})

			It("should skip non-activated available command", func() {
				command.EXPECT().GetActivated().Return(false).AnyTimes()

				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should treat non-executable file as unavailable", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Deactivate(command.GetName())
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion())

				Expect(os.Chmod(command.GetLocation(), 0644)).To(Succeed())
				Expect(doctor.Fix(false)).To(Succeed())
			})

			It("should treat directory as unavailable", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Deactivate(command.GetName())
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion())

				Expect(os.Remove(command.GetLocation())).To(Succeed())
				Expect(os.Mkdir(command.GetLocation(), 0755)).To(Succeed())
				Expect(doctor.Fix(false)).To(Succeed())
			})
		})

		Context("Backup", func() {
			var backupRootDir string

			BeforeEach(func() {
				backupRootDir, _ = os.MkdirTemp("", "cmdr-backup-test")
				Expect(os.WriteFile(filepath.Join(backupRootDir, "testfile"), []byte("test"), 0644)).To(Succeed())
				doctor = manager.NewCommandDoctor(mgr, backupRootDir)
			})

			AfterEach(func() {
				matches, _ := filepath.Glob(backupRootDir + ".backup.*")
				for _, m := range matches {
					os.RemoveAll(m)
				}
				os.RemoveAll(backupRootDir)
			})

			It("should create backup before fix", func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{}, nil)

				Expect(doctor.FixWithOptions(false, true)).To(Succeed())

				matches, err := filepath.Glob(backupRootDir + ".backup.*")
				Expect(err).NotTo(HaveOccurred())
				Expect(matches).To(HaveLen(1))

				backupFile := filepath.Join(matches[0], "testfile")
				Expect(backupFile).To(BeAnExistingFile())
			})

			It("should skip backup when no-backup is true", func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{}, nil)

				Expect(doctor.FixWithOptions(false, false)).To(Succeed())

				matches, _ := filepath.Glob(backupRootDir + ".backup.*")
				Expect(matches).To(BeEmpty())
			})

			It("should skip backup in dry-run mode", func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{}, nil)

				Expect(doctor.FixWithOptions(true, true)).To(Succeed())

				matches, _ := filepath.Glob(backupRootDir + ".backup.*")
				Expect(matches).To(BeEmpty())
			})

			It("should fail when backup root is a file", func() {
				fileRoot := filepath.Join(rootDir, "not-dir")
				Expect(os.WriteFile(fileRoot, []byte("x"), 0644)).To(Succeed())
				doctor = manager.NewCommandDoctor(mgr, fileRoot)

				Expect(doctor.FixWithOptions(false, true)).To(HaveOccurred())
			})

			It("should ignore missing backup roots", func() {
				missingRoot := filepath.Join(rootDir, "missing")
				doctor = manager.NewCommandDoctor(mgr, missingRoot)
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{}, nil)

				Expect(doctor.FixWithOptions(false, true)).To(Succeed())
			})
		})
	})
})
