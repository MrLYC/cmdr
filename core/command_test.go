package core_test

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

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
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

	setContext := func(keys ...define.ContextKey) {
		for _, key := range keys {
			switch key {
			case define.ContextKeyName:
				ctx = context.WithValue(ctx, define.ContextKeyName, name)
			case define.ContextKeyVersion:
				ctx = context.WithValue(ctx, define.ContextKeyVersion, version)
			case define.ContextKeyLocation:
				ctx = context.WithValue(ctx, define.ContextKeyLocation, location)
			case define.ContextKeyDBClient:
				ctx = context.WithValue(ctx, define.ContextKeyDBClient, db)
			case define.ContextKeyCommands:
				ctx = context.WithValue(ctx, define.ContextKeyCommands, commands)
			}
		}
	}

	Context("CommandDefiner", func() {
		var (
			definer *core.CommandDefiner
		)

		BeforeEach(func() {
			definer = core.NewCommandDefiner(shimsDir)
			setContext(define.ContextKeyName, define.ContextKeyVersion, define.ContextKeyLocation, define.ContextKeyDBClient)
		})

		It("name not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyName, "")

			_, err := definer.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
		})

		It("version not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyName, "")

			_, err := definer.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
		})

		It("should define managed command", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommandManaged, true)
			dbQuery.EXPECT().First(gomock.Any()).Return(nil)
			db.EXPECT().Save(gomock.Any()).Return(nil)

			resultCtx, err := definer.Run(ctx)
			Expect(err).To(BeNil())

			commands := resultCtx.Value(define.ContextKeyCommands).([]*model.Command)
			command := commands[0]
			Expect(command.Name).To(Equal(name))
			Expect(command.Version).To(Equal(version))
			Expect(command.Location).To(Equal(filepath.Join(shimsDir, name, fmt.Sprintf("%s_%s", name, version))))
		})

		It("should define unmanaged command", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommandManaged, false)
			dbQuery.EXPECT().First(gomock.Any()).Return(nil)
			db.EXPECT().Save(gomock.Any()).Return(nil)

			resultCtx, err := definer.Run(ctx)
			Expect(err).To(BeNil())

			commands := resultCtx.Value(define.ContextKeyCommands).([]*model.Command)
			command := commands[0]
			Expect(command.Name).To(Equal(name))
			Expect(command.Version).To(Equal(version))
			Expect(command.Location).To(Equal(location))
		})

		It("query failed", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommandManaged, false)
			dbQuery.EXPECT().First(gomock.Any()).Return(fmt.Errorf("error"))

			_, err := definer.Run(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("should update command", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommandManaged, true)
			dbQuery.EXPECT().First(gomock.Any()).DoAndReturn(func(c interface{}) error {
				command := c.(*model.Command)
				command.ID = 1
				return nil
			})

			db.EXPECT().Save(gomock.Any()).DoAndReturn(func(c interface{}) error {
				command := c.(*model.Command)
				Expect(command.ID).To(Equal(1))
				Expect(command.Name).To(Equal(name))
				Expect(command.Version).To(Equal(version))
				Expect(command.Location).To(Equal(filepath.Join(shimsDir, name, fmt.Sprintf("%s_%s", name, version))))
				Expect(command.Managed).To(BeTrue())

				return nil
			})

			_, err := definer.Run(ctx)
			Expect(err).To(BeNil())
		})
	})

	Context("CommandUndefiner", func() {
		var (
			undefiner *core.CommandUndefiner
		)

		BeforeEach(func() {
			undefiner = core.NewCommandUndefiner()
			setContext(define.ContextKeyDBClient, define.ContextKeyCommands)
		})

		It("commands not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommands, []*model.Command{})
			_, err := undefiner.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
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
			activator *core.CommandActivator
		)

		BeforeEach(func() {
			activator = core.NewCommandActivator()
			setContext(define.ContextKeyDBClient, define.ContextKeyCommands)
		})

		It("commands not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyCommands, []*model.Command{})
			_, err := activator.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
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
			deactivator *core.CommandDeactivator
		)

		BeforeEach(func() {
			deactivator = core.NewCommandDeactivator()
			setContext(define.ContextKeyDBClient, define.ContextKeyName)
		})

		It("name not found", func() {
			ctx = context.WithValue(ctx, define.ContextKeyName, "")

			_, err := deactivator.Run(ctx)
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
		})

		It("commands not found", func() {
			dbQuery.EXPECT().Find(gomock.Any()).Return(storm.ErrNotFound)

			_, err := deactivator.Run(ctx)
			Expect(err).To(BeNil())
		})

		It("query failed", func() {
			dbQuery.EXPECT().Find(gomock.Any()).Return(fmt.Errorf("test"))

			_, err := deactivator.Run(ctx)
			Expect(err).NotTo(BeNil())
		})

		It("deactivated", func() {
			command1.Activated = true
			command2.Activated = true
			histories := make([]int, 0, 2)

			dbQuery.EXPECT().Find(gomock.Any()).DoAndReturn(func(c interface{}) error {
				commandsPtr := c.(*[]*model.Command)
				*commandsPtr = commands

				return nil
			})
			db.EXPECT().Save(gomock.Any()).DoAndReturn(func(data interface{}) error {
				command := data.(*model.Command)
				histories = append(histories, command.ID)
				Expect(command.Activated).To(BeFalse())
				return nil
			}).AnyTimes()

			_, err := deactivator.Run(ctx)
			Expect(err).To(BeNil())
			Expect(histories).To(Equal([]int{command1.ID, command2.ID}))
		})

		It("save failed", func() {
			dbQuery.EXPECT().Find(gomock.Any()).DoAndReturn(func(c interface{}) error {
				commandsPtr := c.(*[]*model.Command)
				*commandsPtr = commands

				return nil
			})

			db.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("test")).AnyTimes()
			_, err := deactivator.Run(ctx)
			Expect(err).NotTo(BeNil())
		})
	})
})
