package manager_test

import (
	"fmt"
	"io/ioutil"
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
				command.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
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
				Expect(result.GetVersion()).To(Equal("1.0.0"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
			})

			It("should return commands from recorder query when main query failed", func() {
				command := mock.NewMockCommand(ctrl)
				command.EXPECT().GetName().Return("command").AnyTimes()
				command.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
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
				Expect(result.GetVersion()).To(Equal("1.0.0"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
			})

			It("should merge commands", func() {
				mainCommand := mock.NewMockCommand(ctrl)
				mainCommand.EXPECT().GetName().Return("command").AnyTimes()
				mainCommand.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
				mainCommand.EXPECT().GetLocation().Return("/path/to/command").AnyTimes()
				mainCommand.EXPECT().GetActivated().Return(false).AnyTimes()

				recorderCommand := mock.NewMockCommand(ctrl)
				recorderCommand.EXPECT().GetName().Return("command").AnyTimes()
				recorderCommand.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
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
				Expect(result.GetVersion()).To(Equal("1.0.0"))
				Expect(result.GetLocation()).To(Equal("/path/to/command"))
				Expect(result.GetActivated()).To(BeTrue())
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
			doctor = manager.NewCommandDoctor(mgr)

			rootDir, _ = ioutil.TempDir("", "")

			command.EXPECT().GetName().Return("command").AnyTimes()
			command.EXPECT().GetVersion().Return("1.0.0").AnyTimes()
			command.EXPECT().GetLocation().Return(filepath.Join(rootDir, "command")).AnyTimes()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
		})

		It("should return error when query failed", func() {
			mgr.EXPECT().Query().Return(nil, fmt.Errorf("error"))

			Expect(doctor.Fix()).NotTo(BeNil())
		})

		It("should return error when get commands failed", func() {
			mgr.EXPECT().Query().Return(query, nil)
			query.EXPECT().All().Return(nil, fmt.Errorf("error"))

			Expect(doctor.Fix()).NotTo(BeNil())
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

				Expect(doctor.Fix()).To(Succeed())
			})

			It("should remove non-activate command", func() {
				command.EXPECT().GetActivated().Return(false).AnyTimes()
				mgr.EXPECT().Undefine(command.GetName(), command.GetVersion())

				Expect(doctor.Fix()).To(Succeed())
			})
		})

		Context("Command available", func() {
			BeforeEach(func() {
				mgr.EXPECT().Query().Return(query, nil)
				query.EXPECT().All().Return([]core.Command{command}, nil)
				Expect(ioutil.WriteFile(command.GetLocation(), []byte("command"), 0644)).To(Succeed())
			})

			It("should re-define activated command", func() {
				command.EXPECT().GetActivated().Return(true).AnyTimes()
				mgr.EXPECT().Define(command.GetName(), command.GetVersion(), command.GetLocation())
				mgr.EXPECT().Activate(command.GetName(), command.GetVersion())

				Expect(doctor.Fix()).To(Succeed())
			})

			It("should re-define non-activate command", func() {
				command.EXPECT().GetActivated().Return(false).AnyTimes()
				mgr.EXPECT().Define(command.GetName(), command.GetVersion(), command.GetLocation())

				Expect(doctor.Fix()).To(Succeed())
			})
		})
	})
})
