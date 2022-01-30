package operator_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/operator/mock"
)

var _ = Describe("Command", func() {
	var (
		ctrl               *gomock.Controller
		db                 *mock.MockDBClient
		dbQuery            *mock.MockQuery
		ctx                context.Context
		command1, command2 *model.Command
		commands           []*model.Command
		shimsDir           string
		location           string
		name               string
		version            string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDBClient(ctrl)
		dbQuery = mock.NewMockQuery(ctrl)
		db.EXPECT().Select(gomock.Any()).Return(dbQuery).AnyTimes()

		tempDir, err := afero.TempDir(define.FS, "", "")
		Expect(err).To(BeNil())
		shimsDir = filepath.Join(tempDir, "shims")
		location = filepath.Join(tempDir, "run.sh")

		ctx = context.Background()
		ctx = context.WithValue(ctx, define.ContextKeyDBClient, db)
		ctx = context.WithValue(ctx, define.ContextKeyCommands, commands)

		name = "test"
		version = "1.0.0"

		location1 := filepath.Join(tempDir, "test1.sh")
		err = afero.WriteFile(define.FS, location1, []byte(`#!/bin/sh\necho $@`), 0755)
		Expect(err).To(BeNil())

		command1 = &model.Command{
			ID:       1,
			Name:     "test1",
			Version:  "1.0.0",
			Location: location1,
			Managed:  true,
		}

		location2 := filepath.Join(tempDir, "test2.sh")
		err = afero.WriteFile(define.FS, location2, []byte(`#!/bin/sh\necho $@`), 0755)
		Expect(err).To(BeNil())

		command2 = &model.Command{
			ID:       2,
			Name:     "test2",
			Version:  "1.0.0",
			Location: location2,
			Managed:  true,
		}

		commands = []*model.Command{command1, command2}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CommandDefiner", func() {
		It("should define managed command", func() {
			dbQuery.EXPECT().First(gomock.Any()).Return(nil)

			definer := operator.NewCommandDefiner(shimsDir, name, version, location, true)
			resultCtx, err := definer.Run(ctx)
			Expect(err).To(BeNil())

			commands := resultCtx.Value(define.ContextKeyCommands).([]*model.Command)
			command := commands[0]
			Expect(command.Name).To(Equal(name))
			Expect(command.Version).To(Equal(version))
			Expect(command.Location).To(Equal(location))
		})

		It("should define unmanaged command", func() {
			dbQuery.EXPECT().First(gomock.Any()).Return(nil)

			definer := operator.NewCommandDefiner(shimsDir, name, version, location, false)
			resultCtx, err := definer.Run(ctx)
			Expect(err).To(BeNil())

			commands := resultCtx.Value(define.ContextKeyCommands).([]*model.Command)
			command := commands[0]
			Expect(command.Name).To(Equal(name))
			Expect(command.Version).To(Equal(version))
			Expect(command.Location).To(Equal(location))
		})

		It("query failed", func() {
			dbQuery.EXPECT().First(gomock.Any()).Return(fmt.Errorf("error"))

			definer := operator.NewCommandDefiner(shimsDir, name, version, location, true)
			_, err := definer.Run(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("should update command", func() {
			db.EXPECT().Save(gomock.Any()).DoAndReturn(func(c interface{}) error {
				command := c.(*model.Command)
				Expect(command.Name).To(Equal(name))
				Expect(command.Version).To(Equal(version))
				Expect(command.Location).To(Equal(filepath.Join(shimsDir, name, fmt.Sprintf("%s_%s", name, version))))
				Expect(command.Managed).To(BeTrue())

				return nil
			})

			definer := operator.NewCommandDefiner(shimsDir, name, version, location, true)
			err := definer.Commit(ctx)
			Expect(err).To(BeNil())
		})
	})

	Context("CommandUndefiner", func() {
		var (
			undefiner *operator.CommandUndefiner
		)

		BeforeEach(func() {
			undefiner = operator.NewCommandUndefiner()
		})

		It("commands not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommands, nil)
			_, err := undefiner.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
		})

		It("commands deleted", func() {
			histories := make([]int, 0, 2)
			db.EXPECT().DeleteStruct(gomock.Any()).DoAndReturn(func(data interface{}) error {
				command := data.(*model.Command)
				histories = append(histories, command.ID)

				return nil
			}).AnyTimes()

			_, err := undefiner.Run(ctx)
			Expect(err).To(BeNil())
			Expect(histories).To(Equal([]int{command1.ID, command2.ID}))
		})

		It("commands not found", func() {
			db.EXPECT().DeleteStruct(gomock.Any()).Return(storm.ErrNotFound).AnyTimes()

			_, err := undefiner.Run(ctx)
			Expect(err).To(BeNil())
		})

		It("run failed", func() {
			db.EXPECT().DeleteStruct(gomock.Any()).Return(fmt.Errorf("test")).AnyTimes()

			_, err := undefiner.Run(ctx)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("CommandActivator", func() {
		var (
			activator *operator.CommandActivator
		)

		BeforeEach(func() {
			activator = operator.NewCommandActivator()
		})

		It("commands not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommands, nil)
			_, err := activator.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
		})

		It("commands activated", func() {
			histories := make([]int, 0, 2)

			db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
				command := data.(*model.Command)
				histories = append(histories, command.ID)
				Expect(command.Activated).To(BeTrue())
				return nil
			}).AnyTimes()

			_, err := activator.Run(ctx)
			Expect(err).To(BeNil())

			Expect(histories).To(Equal([]int{command1.ID, command2.ID}))
		})

		It("run failed", func() {
			db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("test")).AnyTimes()

			_, err := activator.Run(ctx)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("CommandDeactivator", func() {
		var (
			deactivator *operator.CommandsDeactivator
		)

		BeforeEach(func() {
			deactivator = operator.NewCommandDeactivator()
		})

		It("context commands not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommands, nil)

			_, err := deactivator.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
		})

		It("query not found", func() {
			dbQuery.EXPECT().Find(gomock.Any()).Return(storm.ErrNotFound).AnyTimes()

			_, err := deactivator.Run(ctx)
			Expect(err).To(BeNil())
		})

		It("query failed", func() {
			dbQuery.EXPECT().Find(gomock.Any()).Return(fmt.Errorf("test")).AnyTimes()

			_, err := deactivator.Run(ctx)
			Expect(err).NotTo(BeNil())
		})

		It("deactivated", func() {
			count := 0

			dbQuery.EXPECT().Find(gomock.Any()).DoAndReturn(func(c interface{}) error {
				commandsPtr := c.(*[]*model.Command)
				*commandsPtr = []*model.Command{{Activated: true}}

				return nil
			}).AnyTimes()
			db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
				command := data.(*model.Command)
				Expect(command.Activated).To(BeFalse())
				count++
				return nil
			}).AnyTimes()

			_, err := deactivator.Run(ctx)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(2))
		})

		It("save failed", func() {
			dbQuery.EXPECT().Find(gomock.Any()).DoAndReturn(func(c interface{}) error {
				commandsPtr := c.(*[]*model.Command)
				*commandsPtr = []*model.Command{{Activated: true}}

				return nil
			}).AnyTimes()

			db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("test")).AnyTimes()
			_, err := deactivator.Run(ctx)
			Expect(err).NotTo(BeNil())
		})
	})
})
